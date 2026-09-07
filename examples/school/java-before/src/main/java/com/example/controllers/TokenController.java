package com.example.controllers;

import javax.servlet.http.HttpServletRequest;

import com.example.util.JwtUtil;
import com.example.util.RemoteAddressUtil;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/token")
public class TokenController {

  @GetMapping("/")
  public String get(HttpServletRequest request) {
    String remoteAddress = RemoteAddressUtil.get(request);
    return JwtUtil.generateToken(remoteAddress);
  }
}