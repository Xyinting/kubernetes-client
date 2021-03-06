/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package validation

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

type RuntimeObjectValidator interface {
	Validate(obj runtime.Object) field.ErrorList
	ValidateUpdate(obj, old runtime.Object) field.ErrorList
}

var Validator = &RuntimeObjectsValidator{map[reflect.Type]RuntimeObjectValidatorInfo{}}

type RuntimeObjectsValidator struct {
	typeToValidator map[reflect.Type]RuntimeObjectValidatorInfo
}

type RuntimeObjectValidatorInfo struct {
	Validator     RuntimeObjectValidator
	IsNamespaced  bool
	HasObjectMeta bool
	UpdateAllowed bool
}

func (v *RuntimeObjectsValidator) GetInfo(obj runtime.Object) (RuntimeObjectValidatorInfo, bool) {
	ret, ok := v.typeToValidator[reflect.TypeOf(obj)]
	return ret, ok
}

func (v *RuntimeObjectsValidator) MustRegister(obj runtime.Object, validateFunction interface{}, validateUpdateFunction interface{}) {
	if err := v.Register(obj, validateFunction, validateUpdateFunction); err != nil {
		panic(err)
	}
}

func (v *RuntimeObjectsValidator) Register(obj runtime.Object, validateFunction interface{}, validateUpdateFunction interface{}) error {
	objType := reflect.TypeOf(obj)
	if oldValidator, exists := v.typeToValidator[objType]; exists {
		panic(fmt.Sprintf("%v is already registered with %v", objType, oldValidator))
	}

	validator, err := NewValidationWrapper(validateFunction, validateUpdateFunction)
	if err != nil {
		return err
	}

	isNamespaced, err := GetRequiresNamespace(obj)
	if err != nil {
		return err
	}

	updateAllowed := validateUpdateFunction != nil

	v.typeToValidator[objType] = RuntimeObjectValidatorInfo{validator, isNamespaced, HasObjectMeta(obj), updateAllowed}

	return nil
}

func (v *RuntimeObjectsValidator) Validate(obj runtime.Object) field.ErrorList {
	if obj == nil {
		return field.ErrorList{}
	}

	allErrs := field.ErrorList{}

	specificValidationInfo, err := v.getSpecificValidationInfo(obj)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(nil, err))
		return allErrs
	}

	allErrs = append(allErrs, specificValidationInfo.Validator.Validate(obj)...)
	return allErrs
}

func (v *RuntimeObjectsValidator) ValidateUpdate(obj, old runtime.Object) field.ErrorList {
	if obj == nil && old == nil {
		return field.ErrorList{}
	}
	if newType, oldType := reflect.TypeOf(obj), reflect.TypeOf(old); newType != oldType {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), newType.Kind(), fmt.Sprintf("expected type %s, for field %s, got %s", oldType.Kind().String(), "kind", newType.Kind().String()))}
	}

	allErrs := field.ErrorList{}

	specificValidationInfo, err := v.getSpecificValidationInfo(obj)
	if err != nil {
		if fieldErr, ok := err.(*field.Error); ok {
			allErrs = append(allErrs, fieldErr)
		} else {
			allErrs = append(allErrs, field.InternalError(nil, err))
		}
		return allErrs
	}

	allErrs = append(allErrs, specificValidationInfo.Validator.ValidateUpdate(obj, old)...)

	// no errors so far, make sure that the new object is actually valid against the original validator
	if len(allErrs) == 0 {
		allErrs = append(allErrs, specificValidationInfo.Validator.Validate(obj)...)
	}

	return allErrs
}

func (v *RuntimeObjectsValidator) getSpecificValidationInfo(obj runtime.Object) (RuntimeObjectValidatorInfo, error) {
	objType := reflect.TypeOf(obj)
	specificValidationInfo, exists := v.typeToValidator[objType]

	if !exists {
		return RuntimeObjectValidatorInfo{}, fmt.Errorf("no validator registered for %v", objType)
	}

	return specificValidationInfo, nil
}

func GetRequiresNamespace(obj runtime.Object) (bool, error) {
	groupVersionKinds, _, err := legacyscheme.Scheme.ObjectKinds(obj)
	if err != nil {
		return false, err
	}

	for _, gvk := range groupVersionKinds {
		restMapping, err := legacyscheme.Registry.RESTMapper().RESTMapping(gvk.GroupKind())
		if err != nil {
			return false, err
		}
		if restMapping.Scope.Name() == meta.RESTScopeNameNamespace {
			return true, nil
		}
	}

	return false, nil
}

func HasObjectMeta(obj runtime.Object) bool {
	objValue := reflect.ValueOf(obj).Elem()
	return objValue.FieldByName("ObjectMeta").IsValid()
}
