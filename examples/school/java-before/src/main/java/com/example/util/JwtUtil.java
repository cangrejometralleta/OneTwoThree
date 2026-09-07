package com.example.util;

import java.time.Instant;
import java.util.Date;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;

public class JwtUtil {
  private final static String JWT_SECRET = System.getenv().getOrDefault("JWT_SECRET",
      "56a7010456b474aeee111f3b7336581fb0a99129d426cf51903efbdfd629f008");
  private final static long JWT_EXPIRATION = 600_000;

  public static Claims getClaims(String token) {
    return Jwts.parser().setSigningKey(JWT_SECRET).parseClaimsJws(token).getBody();
  }

  public static boolean tokenIsValid(String token, String remoteAddress) {
    Claims claims = getClaims(token);
    String subject = claims.getSubject();
    Instant expiration = claims.getExpiration().toInstant();

    return remoteAddress.equals(subject) && Instant.now().compareTo(expiration) > 0;
  }

  public static String generateToken(String remoteAddress) {
    return Jwts.builder().setSubject(remoteAddress).setIssuedAt(new Date(System.currentTimeMillis()))
        .setExpiration(new Date(System.currentTimeMillis() + JWT_EXPIRATION))
        .signWith(SignatureAlgorithm.HS512, JWT_SECRET).compact();
  }

}