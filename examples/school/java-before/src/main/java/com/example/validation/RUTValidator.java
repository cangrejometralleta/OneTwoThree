
package com.example.validation;

import com.example.util.RUTUtil;

import javax.validation.ConstraintValidator;
import javax.validation.ConstraintValidatorContext;

public class RUTValidator implements ConstraintValidator<RUTConstraint, String> {

  @Override
  public void initialize(RUTConstraint rut) {
  }

  @Override
  public boolean isValid(String rut, ConstraintValidatorContext cxt) {
    return RUTUtil.validate(rut);
  }

}