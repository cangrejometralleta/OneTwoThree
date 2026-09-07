package com.example.util;

import javax.servlet.http.HttpServletRequest;

public class RemoteAddressUtil {

  public static String get(HttpServletRequest request) {
    String remoteAddr = "";
    if (request != null) {
      remoteAddr = request.getHeader("X-FORWARDED-FOR");
      if (remoteAddr == null || "".equals(remoteAddr)) {
        remoteAddr = request.getRemoteAddr();
      }
    }
    return remoteAddr;
  }
}