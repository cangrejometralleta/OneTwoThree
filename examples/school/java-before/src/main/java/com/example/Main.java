package com.example;

import com.example.security.JwtAuthorizationFilter;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.builders.WebSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

@SpringBootApplication
public class Main {

  public static void main(String[] args) {
    SpringApplication.run(Main.class, args);
  }

  @EnableWebSecurity
  @Configuration
  class WebSecurityConfig extends WebSecurityConfigurerAdapter {
    @Override
    protected void configure(HttpSecurity http) throws Exception {
      http.antMatcher("*/*").addFilterAfter(new JwtAuthorizationFilter(), UsernamePasswordAuthenticationFilter.class)
          .authorizeRequests().antMatchers("*/*").authenticated();
    }

    @Override
    public void configure(WebSecurity webSecurity) throws Exception {
      webSecurity.ignoring().antMatchers(HttpMethod.GET, "/token/");
    }
  }
}