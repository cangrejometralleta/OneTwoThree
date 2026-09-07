package com.example.security;

import java.io.IOException;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.example.exceptions.TokenNotProvidedException;
import com.example.util.JwtUtil;
import com.example.util.RemoteAddressUtil;

import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.filter.OncePerRequestFilter;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.MalformedJwtException;
import io.jsonwebtoken.UnsupportedJwtException;

public class JwtAuthorizationFilter extends OncePerRequestFilter {

  private final String HEADER = "Authorization";
  private final String PREFIX = "Bearer ";

  @Override
  protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain chain)
      throws ServletException, IOException {
    try {
      if (!"/token/".equals(request.getRequestURI())) {
        String token = getTokenFromHeader(request);
        String remoteAddress = RemoteAddressUtil.get(request);

        if (JwtUtil.tokenIsValid(token, remoteAddress)) {
          Claims claims = JwtUtil.getClaims(token);

          UsernamePasswordAuthenticationToken auth = new UsernamePasswordAuthenticationToken(claims.getSubject(), null);
          SecurityContextHolder.getContext().setAuthentication(auth);
        } else {
          SecurityContextHolder.clearContext();
        }
      }
      chain.doFilter(request, response);
    } catch (ExpiredJwtException | UnsupportedJwtException | MalformedJwtException | TokenNotProvidedException e) {
      response.setStatus(HttpServletResponse.SC_FORBIDDEN);
      ((HttpServletResponse) response).sendError(HttpServletResponse.SC_FORBIDDEN, e.getMessage());
      return;
    }
  }

  private String getTokenFromHeader(HttpServletRequest req) throws TokenNotProvidedException {
    String authenticationHeader = req.getHeader(HEADER);
    if (authenticationHeader == null || !authenticationHeader.startsWith(PREFIX))
      throw new TokenNotProvidedException();
    return authenticationHeader.replace(PREFIX, "");
  }
}