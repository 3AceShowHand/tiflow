/**
 * Autogenerated by Avro
 *
 * DO NOT EDIT DIRECTLY
 */
package com.pingcap.simple.avro;
@org.apache.avro.specific.AvroGenerated
public enum DDLType implements org.apache.avro.generic.GenericEnumSymbol<DDLType> {
  CREATE, ALTER, ERASE, RENAME, TRUNCATE, CINDEX, DINDEX, QUERY  ;
  public static final org.apache.avro.Schema SCHEMA$ = new org.apache.avro.Schema.Parser().parse("{\"type\":\"enum\",\"name\":\"DDLType\",\"namespace\":\"com.pingcap.simple.avro\",\"symbols\":[\"CREATE\",\"ALTER\",\"ERASE\",\"RENAME\",\"TRUNCATE\",\"CINDEX\",\"DINDEX\",\"QUERY\"]}");
  public static org.apache.avro.Schema getClassSchema() { return SCHEMA$; }

  @Override
  public org.apache.avro.Schema getSchema() { return SCHEMA$; }
}
