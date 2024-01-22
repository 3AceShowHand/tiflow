/**
 * Autogenerated by Avro
 *
 * DO NOT EDIT DIRECTLY
 */
package com.pingcap.simple.avro;

import org.apache.avro.generic.GenericArray;
import org.apache.avro.specific.SpecificData;
import org.apache.avro.util.Utf8;
import org.apache.avro.message.BinaryMessageEncoder;
import org.apache.avro.message.BinaryMessageDecoder;
import org.apache.avro.message.SchemaStore;

@org.apache.avro.specific.AvroGenerated
public class DML extends org.apache.avro.specific.SpecificRecordBase implements org.apache.avro.specific.SpecificRecord {
  private static final long serialVersionUID = -4117959977261131971L;


  public static final org.apache.avro.Schema SCHEMA$ = new org.apache.avro.Schema.Parser().parse("{\"type\":\"record\",\"name\":\"DML\",\"namespace\":\"com.pingcap.simple.avro\",\"fields\":[{\"name\":\"version\",\"type\":\"int\"},{\"name\":\"database\",\"type\":\"string\"},{\"name\":\"table\",\"type\":\"string\"},{\"name\":\"tableID\",\"type\":\"long\"},{\"name\":\"type\",\"type\":{\"type\":\"enum\",\"name\":\"DMLType\",\"symbols\":[\"INSERT\",\"UPDATE\",\"DELETE\"]}},{\"name\":\"commitTs\",\"type\":\"long\"},{\"name\":\"buildTs\",\"type\":\"long\"},{\"name\":\"schemaVersion\",\"type\":\"long\"},{\"name\":\"claimCheckLocation\",\"type\":[\"null\",\"string\"],\"default\":null},{\"name\":\"handleKeyOnly\",\"type\":[\"null\",\"boolean\"],\"default\":null},{\"name\":\"checksum\",\"type\":[\"null\",{\"type\":\"record\",\"name\":\"Checksum\",\"fields\":[{\"name\":\"version\",\"type\":\"int\"},{\"name\":\"corrupted\",\"type\":\"boolean\"},{\"name\":\"current\",\"type\":\"long\"},{\"name\":\"previous\",\"type\":\"long\"}],\"docs\":\"event's e2e checksum information\"}],\"default\":null},{\"name\":\"data\",\"type\":[\"null\",{\"type\":\"map\",\"values\":[\"null\",\"int\",\"long\",\"float\",\"double\",\"string\",\"boolean\",\"bytes\"],\"default\":null}],\"default\":null},{\"name\":\"old\",\"type\":[\"null\",{\"type\":\"map\",\"values\":[\"null\",\"int\",\"long\",\"float\",\"double\",\"string\",\"boolean\",\"bytes\"],\"default\":null}],\"default\":null}],\"docs\":\"the message format of the DML event\"}");
  public static org.apache.avro.Schema getClassSchema() { return SCHEMA$; }

  private static final SpecificData MODEL$ = new SpecificData();

  private static final BinaryMessageEncoder<DML> ENCODER =
      new BinaryMessageEncoder<>(MODEL$, SCHEMA$);

  private static final BinaryMessageDecoder<DML> DECODER =
      new BinaryMessageDecoder<>(MODEL$, SCHEMA$);

  /**
   * Return the BinaryMessageEncoder instance used by this class.
   * @return the message encoder used by this class
   */
  public static BinaryMessageEncoder<DML> getEncoder() {
    return ENCODER;
  }

  /**
   * Return the BinaryMessageDecoder instance used by this class.
   * @return the message decoder used by this class
   */
  public static BinaryMessageDecoder<DML> getDecoder() {
    return DECODER;
  }

  /**
   * Create a new BinaryMessageDecoder instance for this class that uses the specified {@link SchemaStore}.
   * @param resolver a {@link SchemaStore} used to find schemas by fingerprint
   * @return a BinaryMessageDecoder instance for this class backed by the given SchemaStore
   */
  public static BinaryMessageDecoder<DML> createDecoder(SchemaStore resolver) {
    return new BinaryMessageDecoder<>(MODEL$, SCHEMA$, resolver);
  }

  /**
   * Serializes this DML to a ByteBuffer.
   * @return a buffer holding the serialized data for this instance
   * @throws java.io.IOException if this instance could not be serialized
   */
  public java.nio.ByteBuffer toByteBuffer() throws java.io.IOException {
    return ENCODER.encode(this);
  }

  /**
   * Deserializes a DML from a ByteBuffer.
   * @param b a byte buffer holding serialized data for an instance of this class
   * @return a DML instance decoded from the given buffer
   * @throws java.io.IOException if the given bytes could not be deserialized into an instance of this class
   */
  public static DML fromByteBuffer(
      java.nio.ByteBuffer b) throws java.io.IOException {
    return DECODER.decode(b);
  }

  private int version;
  private java.lang.CharSequence database;
  private java.lang.CharSequence table;
  private long tableID;
  private com.pingcap.simple.avro.DMLType type;
  private long commitTs;
  private long buildTs;
  private long schemaVersion;
  private java.lang.CharSequence claimCheckLocation;
  private java.lang.Boolean handleKeyOnly;
  private com.pingcap.simple.avro.Checksum checksum;
  private java.util.Map<java.lang.CharSequence,java.lang.Object> data;
  private java.util.Map<java.lang.CharSequence,java.lang.Object> old;

  /**
   * Default constructor.  Note that this does not initialize fields
   * to their default values from the schema.  If that is desired then
   * one should use <code>newBuilder()</code>.
   */
  public DML() {}

  /**
   * All-args constructor.
   * @param version The new value for version
   * @param database The new value for database
   * @param table The new value for table
   * @param tableID The new value for tableID
   * @param type The new value for type
   * @param commitTs The new value for commitTs
   * @param buildTs The new value for buildTs
   * @param schemaVersion The new value for schemaVersion
   * @param claimCheckLocation The new value for claimCheckLocation
   * @param handleKeyOnly The new value for handleKeyOnly
   * @param checksum The new value for checksum
   * @param data The new value for data
   * @param old The new value for old
   */
  public DML(java.lang.Integer version, java.lang.CharSequence database, java.lang.CharSequence table, java.lang.Long tableID, com.pingcap.simple.avro.DMLType type, java.lang.Long commitTs, java.lang.Long buildTs, java.lang.Long schemaVersion, java.lang.CharSequence claimCheckLocation, java.lang.Boolean handleKeyOnly, com.pingcap.simple.avro.Checksum checksum, java.util.Map<java.lang.CharSequence,java.lang.Object> data, java.util.Map<java.lang.CharSequence,java.lang.Object> old) {
    this.version = version;
    this.database = database;
    this.table = table;
    this.tableID = tableID;
    this.type = type;
    this.commitTs = commitTs;
    this.buildTs = buildTs;
    this.schemaVersion = schemaVersion;
    this.claimCheckLocation = claimCheckLocation;
    this.handleKeyOnly = handleKeyOnly;
    this.checksum = checksum;
    this.data = data;
    this.old = old;
  }

  @Override
  public org.apache.avro.specific.SpecificData getSpecificData() { return MODEL$; }

  @Override
  public org.apache.avro.Schema getSchema() { return SCHEMA$; }

  // Used by DatumWriter.  Applications should not call.
  @Override
  public java.lang.Object get(int field$) {
    switch (field$) {
    case 0: return version;
    case 1: return database;
    case 2: return table;
    case 3: return tableID;
    case 4: return type;
    case 5: return commitTs;
    case 6: return buildTs;
    case 7: return schemaVersion;
    case 8: return claimCheckLocation;
    case 9: return handleKeyOnly;
    case 10: return checksum;
    case 11: return data;
    case 12: return old;
    default: throw new IndexOutOfBoundsException("Invalid index: " + field$);
    }
  }

  // Used by DatumReader.  Applications should not call.
  @Override
  @SuppressWarnings(value="unchecked")
  public void put(int field$, java.lang.Object value$) {
    switch (field$) {
    case 0: version = (java.lang.Integer)value$; break;
    case 1: database = (java.lang.CharSequence)value$; break;
    case 2: table = (java.lang.CharSequence)value$; break;
    case 3: tableID = (java.lang.Long)value$; break;
    case 4: type = (com.pingcap.simple.avro.DMLType)value$; break;
    case 5: commitTs = (java.lang.Long)value$; break;
    case 6: buildTs = (java.lang.Long)value$; break;
    case 7: schemaVersion = (java.lang.Long)value$; break;
    case 8: claimCheckLocation = (java.lang.CharSequence)value$; break;
    case 9: handleKeyOnly = (java.lang.Boolean)value$; break;
    case 10: checksum = (com.pingcap.simple.avro.Checksum)value$; break;
    case 11: data = (java.util.Map<java.lang.CharSequence,java.lang.Object>)value$; break;
    case 12: old = (java.util.Map<java.lang.CharSequence,java.lang.Object>)value$; break;
    default: throw new IndexOutOfBoundsException("Invalid index: " + field$);
    }
  }

  /**
   * Gets the value of the 'version' field.
   * @return The value of the 'version' field.
   */
  public int getVersion() {
    return version;
  }


  /**
   * Sets the value of the 'version' field.
   * @param value the value to set.
   */
  public void setVersion(int value) {
    this.version = value;
  }

  /**
   * Gets the value of the 'database' field.
   * @return The value of the 'database' field.
   */
  public java.lang.CharSequence getDatabase() {
    return database;
  }


  /**
   * Sets the value of the 'database' field.
   * @param value the value to set.
   */
  public void setDatabase(java.lang.CharSequence value) {
    this.database = value;
  }

  /**
   * Gets the value of the 'table' field.
   * @return The value of the 'table' field.
   */
  public java.lang.CharSequence getTable() {
    return table;
  }


  /**
   * Sets the value of the 'table' field.
   * @param value the value to set.
   */
  public void setTable(java.lang.CharSequence value) {
    this.table = value;
  }

  /**
   * Gets the value of the 'tableID' field.
   * @return The value of the 'tableID' field.
   */
  public long getTableID() {
    return tableID;
  }


  /**
   * Sets the value of the 'tableID' field.
   * @param value the value to set.
   */
  public void setTableID(long value) {
    this.tableID = value;
  }

  /**
   * Gets the value of the 'type' field.
   * @return The value of the 'type' field.
   */
  public com.pingcap.simple.avro.DMLType getType() {
    return type;
  }


  /**
   * Sets the value of the 'type' field.
   * @param value the value to set.
   */
  public void setType(com.pingcap.simple.avro.DMLType value) {
    this.type = value;
  }

  /**
   * Gets the value of the 'commitTs' field.
   * @return The value of the 'commitTs' field.
   */
  public long getCommitTs() {
    return commitTs;
  }


  /**
   * Sets the value of the 'commitTs' field.
   * @param value the value to set.
   */
  public void setCommitTs(long value) {
    this.commitTs = value;
  }

  /**
   * Gets the value of the 'buildTs' field.
   * @return The value of the 'buildTs' field.
   */
  public long getBuildTs() {
    return buildTs;
  }


  /**
   * Sets the value of the 'buildTs' field.
   * @param value the value to set.
   */
  public void setBuildTs(long value) {
    this.buildTs = value;
  }

  /**
   * Gets the value of the 'schemaVersion' field.
   * @return The value of the 'schemaVersion' field.
   */
  public long getSchemaVersion() {
    return schemaVersion;
  }


  /**
   * Sets the value of the 'schemaVersion' field.
   * @param value the value to set.
   */
  public void setSchemaVersion(long value) {
    this.schemaVersion = value;
  }

  /**
   * Gets the value of the 'claimCheckLocation' field.
   * @return The value of the 'claimCheckLocation' field.
   */
  public java.lang.CharSequence getClaimCheckLocation() {
    return claimCheckLocation;
  }


  /**
   * Sets the value of the 'claimCheckLocation' field.
   * @param value the value to set.
   */
  public void setClaimCheckLocation(java.lang.CharSequence value) {
    this.claimCheckLocation = value;
  }

  /**
   * Gets the value of the 'handleKeyOnly' field.
   * @return The value of the 'handleKeyOnly' field.
   */
  public java.lang.Boolean getHandleKeyOnly() {
    return handleKeyOnly;
  }


  /**
   * Sets the value of the 'handleKeyOnly' field.
   * @param value the value to set.
   */
  public void setHandleKeyOnly(java.lang.Boolean value) {
    this.handleKeyOnly = value;
  }

  /**
   * Gets the value of the 'checksum' field.
   * @return The value of the 'checksum' field.
   */
  public com.pingcap.simple.avro.Checksum getChecksum() {
    return checksum;
  }


  /**
   * Sets the value of the 'checksum' field.
   * @param value the value to set.
   */
  public void setChecksum(com.pingcap.simple.avro.Checksum value) {
    this.checksum = value;
  }

  /**
   * Gets the value of the 'data' field.
   * @return The value of the 'data' field.
   */
  public java.util.Map<java.lang.CharSequence,java.lang.Object> getData() {
    return data;
  }


  /**
   * Sets the value of the 'data' field.
   * @param value the value to set.
   */
  public void setData(java.util.Map<java.lang.CharSequence,java.lang.Object> value) {
    this.data = value;
  }

  /**
   * Gets the value of the 'old' field.
   * @return The value of the 'old' field.
   */
  public java.util.Map<java.lang.CharSequence,java.lang.Object> getOld() {
    return old;
  }


  /**
   * Sets the value of the 'old' field.
   * @param value the value to set.
   */
  public void setOld(java.util.Map<java.lang.CharSequence,java.lang.Object> value) {
    this.old = value;
  }

  /**
   * Creates a new DML RecordBuilder.
   * @return A new DML RecordBuilder
   */
  public static com.pingcap.simple.avro.DML.Builder newBuilder() {
    return new com.pingcap.simple.avro.DML.Builder();
  }

  /**
   * Creates a new DML RecordBuilder by copying an existing Builder.
   * @param other The existing builder to copy.
   * @return A new DML RecordBuilder
   */
  public static com.pingcap.simple.avro.DML.Builder newBuilder(com.pingcap.simple.avro.DML.Builder other) {
    if (other == null) {
      return new com.pingcap.simple.avro.DML.Builder();
    } else {
      return new com.pingcap.simple.avro.DML.Builder(other);
    }
  }

  /**
   * Creates a new DML RecordBuilder by copying an existing DML instance.
   * @param other The existing instance to copy.
   * @return A new DML RecordBuilder
   */
  public static com.pingcap.simple.avro.DML.Builder newBuilder(com.pingcap.simple.avro.DML other) {
    if (other == null) {
      return new com.pingcap.simple.avro.DML.Builder();
    } else {
      return new com.pingcap.simple.avro.DML.Builder(other);
    }
  }

  /**
   * RecordBuilder for DML instances.
   */
  @org.apache.avro.specific.AvroGenerated
  public static class Builder extends org.apache.avro.specific.SpecificRecordBuilderBase<DML>
    implements org.apache.avro.data.RecordBuilder<DML> {

    private int version;
    private java.lang.CharSequence database;
    private java.lang.CharSequence table;
    private long tableID;
    private com.pingcap.simple.avro.DMLType type;
    private long commitTs;
    private long buildTs;
    private long schemaVersion;
    private java.lang.CharSequence claimCheckLocation;
    private java.lang.Boolean handleKeyOnly;
    private com.pingcap.simple.avro.Checksum checksum;
    private com.pingcap.simple.avro.Checksum.Builder checksumBuilder;
    private java.util.Map<java.lang.CharSequence,java.lang.Object> data;
    private java.util.Map<java.lang.CharSequence,java.lang.Object> old;

    /** Creates a new Builder */
    private Builder() {
      super(SCHEMA$, MODEL$);
    }

    /**
     * Creates a Builder by copying an existing Builder.
     * @param other The existing Builder to copy.
     */
    private Builder(com.pingcap.simple.avro.DML.Builder other) {
      super(other);
      if (isValidValue(fields()[0], other.version)) {
        this.version = data().deepCopy(fields()[0].schema(), other.version);
        fieldSetFlags()[0] = other.fieldSetFlags()[0];
      }
      if (isValidValue(fields()[1], other.database)) {
        this.database = data().deepCopy(fields()[1].schema(), other.database);
        fieldSetFlags()[1] = other.fieldSetFlags()[1];
      }
      if (isValidValue(fields()[2], other.table)) {
        this.table = data().deepCopy(fields()[2].schema(), other.table);
        fieldSetFlags()[2] = other.fieldSetFlags()[2];
      }
      if (isValidValue(fields()[3], other.tableID)) {
        this.tableID = data().deepCopy(fields()[3].schema(), other.tableID);
        fieldSetFlags()[3] = other.fieldSetFlags()[3];
      }
      if (isValidValue(fields()[4], other.type)) {
        this.type = data().deepCopy(fields()[4].schema(), other.type);
        fieldSetFlags()[4] = other.fieldSetFlags()[4];
      }
      if (isValidValue(fields()[5], other.commitTs)) {
        this.commitTs = data().deepCopy(fields()[5].schema(), other.commitTs);
        fieldSetFlags()[5] = other.fieldSetFlags()[5];
      }
      if (isValidValue(fields()[6], other.buildTs)) {
        this.buildTs = data().deepCopy(fields()[6].schema(), other.buildTs);
        fieldSetFlags()[6] = other.fieldSetFlags()[6];
      }
      if (isValidValue(fields()[7], other.schemaVersion)) {
        this.schemaVersion = data().deepCopy(fields()[7].schema(), other.schemaVersion);
        fieldSetFlags()[7] = other.fieldSetFlags()[7];
      }
      if (isValidValue(fields()[8], other.claimCheckLocation)) {
        this.claimCheckLocation = data().deepCopy(fields()[8].schema(), other.claimCheckLocation);
        fieldSetFlags()[8] = other.fieldSetFlags()[8];
      }
      if (isValidValue(fields()[9], other.handleKeyOnly)) {
        this.handleKeyOnly = data().deepCopy(fields()[9].schema(), other.handleKeyOnly);
        fieldSetFlags()[9] = other.fieldSetFlags()[9];
      }
      if (isValidValue(fields()[10], other.checksum)) {
        this.checksum = data().deepCopy(fields()[10].schema(), other.checksum);
        fieldSetFlags()[10] = other.fieldSetFlags()[10];
      }
      if (other.hasChecksumBuilder()) {
        this.checksumBuilder = com.pingcap.simple.avro.Checksum.newBuilder(other.getChecksumBuilder());
      }
      if (isValidValue(fields()[11], other.data)) {
        this.data = data().deepCopy(fields()[11].schema(), other.data);
        fieldSetFlags()[11] = other.fieldSetFlags()[11];
      }
      if (isValidValue(fields()[12], other.old)) {
        this.old = data().deepCopy(fields()[12].schema(), other.old);
        fieldSetFlags()[12] = other.fieldSetFlags()[12];
      }
    }

    /**
     * Creates a Builder by copying an existing DML instance
     * @param other The existing instance to copy.
     */
    private Builder(com.pingcap.simple.avro.DML other) {
      super(SCHEMA$, MODEL$);
      if (isValidValue(fields()[0], other.version)) {
        this.version = data().deepCopy(fields()[0].schema(), other.version);
        fieldSetFlags()[0] = true;
      }
      if (isValidValue(fields()[1], other.database)) {
        this.database = data().deepCopy(fields()[1].schema(), other.database);
        fieldSetFlags()[1] = true;
      }
      if (isValidValue(fields()[2], other.table)) {
        this.table = data().deepCopy(fields()[2].schema(), other.table);
        fieldSetFlags()[2] = true;
      }
      if (isValidValue(fields()[3], other.tableID)) {
        this.tableID = data().deepCopy(fields()[3].schema(), other.tableID);
        fieldSetFlags()[3] = true;
      }
      if (isValidValue(fields()[4], other.type)) {
        this.type = data().deepCopy(fields()[4].schema(), other.type);
        fieldSetFlags()[4] = true;
      }
      if (isValidValue(fields()[5], other.commitTs)) {
        this.commitTs = data().deepCopy(fields()[5].schema(), other.commitTs);
        fieldSetFlags()[5] = true;
      }
      if (isValidValue(fields()[6], other.buildTs)) {
        this.buildTs = data().deepCopy(fields()[6].schema(), other.buildTs);
        fieldSetFlags()[6] = true;
      }
      if (isValidValue(fields()[7], other.schemaVersion)) {
        this.schemaVersion = data().deepCopy(fields()[7].schema(), other.schemaVersion);
        fieldSetFlags()[7] = true;
      }
      if (isValidValue(fields()[8], other.claimCheckLocation)) {
        this.claimCheckLocation = data().deepCopy(fields()[8].schema(), other.claimCheckLocation);
        fieldSetFlags()[8] = true;
      }
      if (isValidValue(fields()[9], other.handleKeyOnly)) {
        this.handleKeyOnly = data().deepCopy(fields()[9].schema(), other.handleKeyOnly);
        fieldSetFlags()[9] = true;
      }
      if (isValidValue(fields()[10], other.checksum)) {
        this.checksum = data().deepCopy(fields()[10].schema(), other.checksum);
        fieldSetFlags()[10] = true;
      }
      this.checksumBuilder = null;
      if (isValidValue(fields()[11], other.data)) {
        this.data = data().deepCopy(fields()[11].schema(), other.data);
        fieldSetFlags()[11] = true;
      }
      if (isValidValue(fields()[12], other.old)) {
        this.old = data().deepCopy(fields()[12].schema(), other.old);
        fieldSetFlags()[12] = true;
      }
    }

    /**
      * Gets the value of the 'version' field.
      * @return The value.
      */
    public int getVersion() {
      return version;
    }


    /**
      * Sets the value of the 'version' field.
      * @param value The value of 'version'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setVersion(int value) {
      validate(fields()[0], value);
      this.version = value;
      fieldSetFlags()[0] = true;
      return this;
    }

    /**
      * Checks whether the 'version' field has been set.
      * @return True if the 'version' field has been set, false otherwise.
      */
    public boolean hasVersion() {
      return fieldSetFlags()[0];
    }


    /**
      * Clears the value of the 'version' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearVersion() {
      fieldSetFlags()[0] = false;
      return this;
    }

    /**
      * Gets the value of the 'database' field.
      * @return The value.
      */
    public java.lang.CharSequence getDatabase() {
      return database;
    }


    /**
      * Sets the value of the 'database' field.
      * @param value The value of 'database'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setDatabase(java.lang.CharSequence value) {
      validate(fields()[1], value);
      this.database = value;
      fieldSetFlags()[1] = true;
      return this;
    }

    /**
      * Checks whether the 'database' field has been set.
      * @return True if the 'database' field has been set, false otherwise.
      */
    public boolean hasDatabase() {
      return fieldSetFlags()[1];
    }


    /**
      * Clears the value of the 'database' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearDatabase() {
      database = null;
      fieldSetFlags()[1] = false;
      return this;
    }

    /**
      * Gets the value of the 'table' field.
      * @return The value.
      */
    public java.lang.CharSequence getTable() {
      return table;
    }


    /**
      * Sets the value of the 'table' field.
      * @param value The value of 'table'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setTable(java.lang.CharSequence value) {
      validate(fields()[2], value);
      this.table = value;
      fieldSetFlags()[2] = true;
      return this;
    }

    /**
      * Checks whether the 'table' field has been set.
      * @return True if the 'table' field has been set, false otherwise.
      */
    public boolean hasTable() {
      return fieldSetFlags()[2];
    }


    /**
      * Clears the value of the 'table' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearTable() {
      table = null;
      fieldSetFlags()[2] = false;
      return this;
    }

    /**
      * Gets the value of the 'tableID' field.
      * @return The value.
      */
    public long getTableID() {
      return tableID;
    }


    /**
      * Sets the value of the 'tableID' field.
      * @param value The value of 'tableID'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setTableID(long value) {
      validate(fields()[3], value);
      this.tableID = value;
      fieldSetFlags()[3] = true;
      return this;
    }

    /**
      * Checks whether the 'tableID' field has been set.
      * @return True if the 'tableID' field has been set, false otherwise.
      */
    public boolean hasTableID() {
      return fieldSetFlags()[3];
    }


    /**
      * Clears the value of the 'tableID' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearTableID() {
      fieldSetFlags()[3] = false;
      return this;
    }

    /**
      * Gets the value of the 'type' field.
      * @return The value.
      */
    public com.pingcap.simple.avro.DMLType getType() {
      return type;
    }


    /**
      * Sets the value of the 'type' field.
      * @param value The value of 'type'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setType(com.pingcap.simple.avro.DMLType value) {
      validate(fields()[4], value);
      this.type = value;
      fieldSetFlags()[4] = true;
      return this;
    }

    /**
      * Checks whether the 'type' field has been set.
      * @return True if the 'type' field has been set, false otherwise.
      */
    public boolean hasType() {
      return fieldSetFlags()[4];
    }


    /**
      * Clears the value of the 'type' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearType() {
      type = null;
      fieldSetFlags()[4] = false;
      return this;
    }

    /**
      * Gets the value of the 'commitTs' field.
      * @return The value.
      */
    public long getCommitTs() {
      return commitTs;
    }


    /**
      * Sets the value of the 'commitTs' field.
      * @param value The value of 'commitTs'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setCommitTs(long value) {
      validate(fields()[5], value);
      this.commitTs = value;
      fieldSetFlags()[5] = true;
      return this;
    }

    /**
      * Checks whether the 'commitTs' field has been set.
      * @return True if the 'commitTs' field has been set, false otherwise.
      */
    public boolean hasCommitTs() {
      return fieldSetFlags()[5];
    }


    /**
      * Clears the value of the 'commitTs' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearCommitTs() {
      fieldSetFlags()[5] = false;
      return this;
    }

    /**
      * Gets the value of the 'buildTs' field.
      * @return The value.
      */
    public long getBuildTs() {
      return buildTs;
    }


    /**
      * Sets the value of the 'buildTs' field.
      * @param value The value of 'buildTs'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setBuildTs(long value) {
      validate(fields()[6], value);
      this.buildTs = value;
      fieldSetFlags()[6] = true;
      return this;
    }

    /**
      * Checks whether the 'buildTs' field has been set.
      * @return True if the 'buildTs' field has been set, false otherwise.
      */
    public boolean hasBuildTs() {
      return fieldSetFlags()[6];
    }


    /**
      * Clears the value of the 'buildTs' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearBuildTs() {
      fieldSetFlags()[6] = false;
      return this;
    }

    /**
      * Gets the value of the 'schemaVersion' field.
      * @return The value.
      */
    public long getSchemaVersion() {
      return schemaVersion;
    }


    /**
      * Sets the value of the 'schemaVersion' field.
      * @param value The value of 'schemaVersion'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setSchemaVersion(long value) {
      validate(fields()[7], value);
      this.schemaVersion = value;
      fieldSetFlags()[7] = true;
      return this;
    }

    /**
      * Checks whether the 'schemaVersion' field has been set.
      * @return True if the 'schemaVersion' field has been set, false otherwise.
      */
    public boolean hasSchemaVersion() {
      return fieldSetFlags()[7];
    }


    /**
      * Clears the value of the 'schemaVersion' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearSchemaVersion() {
      fieldSetFlags()[7] = false;
      return this;
    }

    /**
      * Gets the value of the 'claimCheckLocation' field.
      * @return The value.
      */
    public java.lang.CharSequence getClaimCheckLocation() {
      return claimCheckLocation;
    }


    /**
      * Sets the value of the 'claimCheckLocation' field.
      * @param value The value of 'claimCheckLocation'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setClaimCheckLocation(java.lang.CharSequence value) {
      validate(fields()[8], value);
      this.claimCheckLocation = value;
      fieldSetFlags()[8] = true;
      return this;
    }

    /**
      * Checks whether the 'claimCheckLocation' field has been set.
      * @return True if the 'claimCheckLocation' field has been set, false otherwise.
      */
    public boolean hasClaimCheckLocation() {
      return fieldSetFlags()[8];
    }


    /**
      * Clears the value of the 'claimCheckLocation' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearClaimCheckLocation() {
      claimCheckLocation = null;
      fieldSetFlags()[8] = false;
      return this;
    }

    /**
      * Gets the value of the 'handleKeyOnly' field.
      * @return The value.
      */
    public java.lang.Boolean getHandleKeyOnly() {
      return handleKeyOnly;
    }


    /**
      * Sets the value of the 'handleKeyOnly' field.
      * @param value The value of 'handleKeyOnly'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setHandleKeyOnly(java.lang.Boolean value) {
      validate(fields()[9], value);
      this.handleKeyOnly = value;
      fieldSetFlags()[9] = true;
      return this;
    }

    /**
      * Checks whether the 'handleKeyOnly' field has been set.
      * @return True if the 'handleKeyOnly' field has been set, false otherwise.
      */
    public boolean hasHandleKeyOnly() {
      return fieldSetFlags()[9];
    }


    /**
      * Clears the value of the 'handleKeyOnly' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearHandleKeyOnly() {
      handleKeyOnly = null;
      fieldSetFlags()[9] = false;
      return this;
    }

    /**
      * Gets the value of the 'checksum' field.
      * @return The value.
      */
    public com.pingcap.simple.avro.Checksum getChecksum() {
      return checksum;
    }


    /**
      * Sets the value of the 'checksum' field.
      * @param value The value of 'checksum'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setChecksum(com.pingcap.simple.avro.Checksum value) {
      validate(fields()[10], value);
      this.checksumBuilder = null;
      this.checksum = value;
      fieldSetFlags()[10] = true;
      return this;
    }

    /**
      * Checks whether the 'checksum' field has been set.
      * @return True if the 'checksum' field has been set, false otherwise.
      */
    public boolean hasChecksum() {
      return fieldSetFlags()[10];
    }

    /**
     * Gets the Builder instance for the 'checksum' field and creates one if it doesn't exist yet.
     * @return This builder.
     */
    public com.pingcap.simple.avro.Checksum.Builder getChecksumBuilder() {
      if (checksumBuilder == null) {
        if (hasChecksum()) {
          setChecksumBuilder(com.pingcap.simple.avro.Checksum.newBuilder(checksum));
        } else {
          setChecksumBuilder(com.pingcap.simple.avro.Checksum.newBuilder());
        }
      }
      return checksumBuilder;
    }

    /**
     * Sets the Builder instance for the 'checksum' field
     * @param value The builder instance that must be set.
     * @return This builder.
     */

    public com.pingcap.simple.avro.DML.Builder setChecksumBuilder(com.pingcap.simple.avro.Checksum.Builder value) {
      clearChecksum();
      checksumBuilder = value;
      return this;
    }

    /**
     * Checks whether the 'checksum' field has an active Builder instance
     * @return True if the 'checksum' field has an active Builder instance
     */
    public boolean hasChecksumBuilder() {
      return checksumBuilder != null;
    }

    /**
      * Clears the value of the 'checksum' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearChecksum() {
      checksum = null;
      checksumBuilder = null;
      fieldSetFlags()[10] = false;
      return this;
    }

    /**
      * Gets the value of the 'data' field.
      * @return The value.
      */
    public java.util.Map<java.lang.CharSequence,java.lang.Object> getData() {
      return data;
    }


    /**
      * Sets the value of the 'data' field.
      * @param value The value of 'data'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setData(java.util.Map<java.lang.CharSequence,java.lang.Object> value) {
      validate(fields()[11], value);
      this.data = value;
      fieldSetFlags()[11] = true;
      return this;
    }

    /**
      * Checks whether the 'data' field has been set.
      * @return True if the 'data' field has been set, false otherwise.
      */
    public boolean hasData() {
      return fieldSetFlags()[11];
    }


    /**
      * Clears the value of the 'data' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearData() {
      data = null;
      fieldSetFlags()[11] = false;
      return this;
    }

    /**
      * Gets the value of the 'old' field.
      * @return The value.
      */
    public java.util.Map<java.lang.CharSequence,java.lang.Object> getOld() {
      return old;
    }


    /**
      * Sets the value of the 'old' field.
      * @param value The value of 'old'.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder setOld(java.util.Map<java.lang.CharSequence,java.lang.Object> value) {
      validate(fields()[12], value);
      this.old = value;
      fieldSetFlags()[12] = true;
      return this;
    }

    /**
      * Checks whether the 'old' field has been set.
      * @return True if the 'old' field has been set, false otherwise.
      */
    public boolean hasOld() {
      return fieldSetFlags()[12];
    }


    /**
      * Clears the value of the 'old' field.
      * @return This builder.
      */
    public com.pingcap.simple.avro.DML.Builder clearOld() {
      old = null;
      fieldSetFlags()[12] = false;
      return this;
    }

    @Override
    @SuppressWarnings("unchecked")
    public DML build() {
      try {
        DML record = new DML();
        record.version = fieldSetFlags()[0] ? this.version : (java.lang.Integer) defaultValue(fields()[0]);
        record.database = fieldSetFlags()[1] ? this.database : (java.lang.CharSequence) defaultValue(fields()[1]);
        record.table = fieldSetFlags()[2] ? this.table : (java.lang.CharSequence) defaultValue(fields()[2]);
        record.tableID = fieldSetFlags()[3] ? this.tableID : (java.lang.Long) defaultValue(fields()[3]);
        record.type = fieldSetFlags()[4] ? this.type : (com.pingcap.simple.avro.DMLType) defaultValue(fields()[4]);
        record.commitTs = fieldSetFlags()[5] ? this.commitTs : (java.lang.Long) defaultValue(fields()[5]);
        record.buildTs = fieldSetFlags()[6] ? this.buildTs : (java.lang.Long) defaultValue(fields()[6]);
        record.schemaVersion = fieldSetFlags()[7] ? this.schemaVersion : (java.lang.Long) defaultValue(fields()[7]);
        record.claimCheckLocation = fieldSetFlags()[8] ? this.claimCheckLocation : (java.lang.CharSequence) defaultValue(fields()[8]);
        record.handleKeyOnly = fieldSetFlags()[9] ? this.handleKeyOnly : (java.lang.Boolean) defaultValue(fields()[9]);
        if (checksumBuilder != null) {
          try {
            record.checksum = this.checksumBuilder.build();
          } catch (org.apache.avro.AvroMissingFieldException e) {
            e.addParentField(record.getSchema().getField("checksum"));
            throw e;
          }
        } else {
          record.checksum = fieldSetFlags()[10] ? this.checksum : (com.pingcap.simple.avro.Checksum) defaultValue(fields()[10]);
        }
        record.data = fieldSetFlags()[11] ? this.data : (java.util.Map<java.lang.CharSequence,java.lang.Object>) defaultValue(fields()[11]);
        record.old = fieldSetFlags()[12] ? this.old : (java.util.Map<java.lang.CharSequence,java.lang.Object>) defaultValue(fields()[12]);
        return record;
      } catch (org.apache.avro.AvroMissingFieldException e) {
        throw e;
      } catch (java.lang.Exception e) {
        throw new org.apache.avro.AvroRuntimeException(e);
      }
    }
  }

  @SuppressWarnings("unchecked")
  private static final org.apache.avro.io.DatumWriter<DML>
    WRITER$ = (org.apache.avro.io.DatumWriter<DML>)MODEL$.createDatumWriter(SCHEMA$);

  @Override public void writeExternal(java.io.ObjectOutput out)
    throws java.io.IOException {
    WRITER$.write(this, SpecificData.getEncoder(out));
  }

  @SuppressWarnings("unchecked")
  private static final org.apache.avro.io.DatumReader<DML>
    READER$ = (org.apache.avro.io.DatumReader<DML>)MODEL$.createDatumReader(SCHEMA$);

  @Override public void readExternal(java.io.ObjectInput in)
    throws java.io.IOException {
    READER$.read(this, SpecificData.getDecoder(in));
  }

}










