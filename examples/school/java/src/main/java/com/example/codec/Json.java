package com.example.codec;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.StreamReadFeature;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;

import java.nio.ByteBuffer;
import java.nio.charset.CharacterCodingException;
import java.nio.charset.CodingErrorAction;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

/**
 * Json Reads and Writes Plain Values, and is the only Type that Imports Jackson.
 * Reading Hands back a Map, a List, a String, a BigInteger for a whole Number,
 * a Double for any other, a Boolean or NULL. Nothing else Crosses.
 * Reference: https://www.rfc-editor.org/rfc/rfc8259
 */
public final class Json {

  /** NULL Stands for a JSON null, so a Map can Tell it from an absent Key. */
  public static final Object NULL = new Object() {
    @Override
    public String toString() {
      return "null";
    }
  };

  private static final ObjectMapper MAPPER = JsonMapper.builder()
      .enable(DeserializationFeature.FAIL_ON_TRAILING_TOKENS)
      .enable(DeserializationFeature.USE_BIG_INTEGER_FOR_INTS)
      .enable(StreamReadFeature.STRICT_DUPLICATE_DETECTION)
      .disable(JsonParser.Feature.ALLOW_COMMENTS)
      .build();

  private Json() {
  }

  /** readJsonText Reads one whole JSON Value from strict UTF-8, or Nothing. */
  public static Optional<Object> readJsonText(byte[] data) {
    try {
      String text = StandardCharsets.UTF_8.newDecoder()
          .onMalformedInput(CodingErrorAction.REPORT)
          .onUnmappableCharacter(CodingErrorAction.REPORT)
          .decode(ByteBuffer.wrap(data)).toString();
      JsonNode tree = MAPPER.readTree(text);

      return tree == null || tree.isMissingNode() ? Optional.empty() : Optional.of(readPlainValue(tree));
    } catch (CharacterCodingException | JsonProcessingException broken) {
      return Optional.empty();
    }
  }

  /** writeJsonText Renders a Plain Value or a Record as compact JSON. */
  public static String writeJsonText(Object value) {
    try {
      return MAPPER.writeValueAsString(value);
    } catch (JsonProcessingException broken) {
      throw new IllegalArgumentException("value cannot be Written as JSON", broken);
    }
  }

  private static Object readPlainValue(JsonNode node) {
    if (node.isObject()) {
      Map<String, Object> fields = new LinkedHashMap<>();
      for (Iterator<Map.Entry<String, JsonNode>> it = node.fields(); it.hasNext();) {
        Map.Entry<String, JsonNode> field = it.next();
        fields.put(field.getKey(), readPlainValue(field.getValue()));
      }

      return fields;
    }
    if (node.isArray()) {
      List<Object> items = new ArrayList<>();
      node.forEach(item -> items.add(readPlainValue(item)));

      return items;
    }
    if (node.isIntegralNumber()) {
      return node.bigIntegerValue();
    }
    if (node.isNumber()) {
      return node.doubleValue();
    }
    if (node.isTextual()) {
      return node.textValue();
    }
    if (node.isBoolean()) {
      return node.booleanValue();
    }

    return NULL;
  }
}
