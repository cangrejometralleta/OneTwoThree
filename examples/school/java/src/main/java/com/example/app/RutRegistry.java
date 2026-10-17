package com.example.app;

import com.example.school.Rut;

/**
 * RutRegistry Confirms whether a Rut Belongs to someone Real.
 * A Rut it Denies Answers false; a Registry that cannot Answer Throws.
 */
public interface RutRegistry {
  boolean confirmRut(Rut rut) throws Exception;
}
