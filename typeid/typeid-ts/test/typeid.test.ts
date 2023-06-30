import { describe, expect, it } from '@jest/globals';

import { typeid, TypeID } from "../src/typeid";
import validJson from "./valid";
import invalidJson from "./invalid";

describe('TypeID', () => {
  describe('constructor', () => {
    it('should create a TypeID object', () => {
      const prefix = "test";
      const suffix = "00041061050r3gg28a1c60t3gf";

      const id = typeid(prefix, suffix);
      expect(id).toBeInstanceOf(TypeID);
      expect(id.getType()).toEqual(prefix);
      expect(id.getSuffix()).toEqual(suffix);
    });

    it('should generate a suffix when none is provided', () => {
      const prefix = "test";

      const id = typeid(prefix);
      expect(id).toBeInstanceOf(TypeID);
      expect(id.getType()).toEqual(prefix);
      expect(id.getSuffix()).toHaveLength(26);
    });

    it('should throw an error if prefix is not lowercase', () => {
      expect(() => {
        typeid("TEST", "00041061050r3gg28a1c60t3gf");
      }).toThrowError("Invalid prefix. Must be at most 63 ascii letters [a-z]");

      expect(() => {
        typeid("  ", "00041061050r3gg28a1c60t3gf");
      }).toThrowError("Invalid prefix. Must be at most 63 ascii letters [a-z]");
    });

    it('should throw an error if suffix length is not 26', () => {
      expect(() => {
        typeid("test", "abc");
      }).toThrowError("Invalid length. Suffix should have 26 characters, got 3");
    });
  });

  describe('toString', () => {
    it('should return a string representation', () => {
      const prefix = "test";
      const suffix = "00041061050r3gg28a1c60t3gf";

      const id = typeid(prefix, suffix);
      expect(id.toString()).toEqual('test_00041061050r3gg28a1c60t3gf');
    });

    it('should return a string representation even without prefix', () => {
      const suffix = "00041061050r3gg28a1c60t3gf";

      const id = typeid("", suffix);
      expect(id.toString()).toEqual(suffix);
    });
  });

  describe('fromString', () => {
    it('should construct TypeID from a string without prefix', () => {
      const str = '00041061050r3gg28a1c60t3gf';
      const tid = TypeID.fromString(str);

      expect(tid.getSuffix()).toBe(str);
      expect(tid.getType()).toBe('');
    });

    it('should construct TypeID from a string with prefix', () => {
      const str = 'prefix_00041061050r3gg28a1c60t3gf';
      const tid = TypeID.fromString(str);

      expect(tid.getSuffix()).toBe('00041061050r3gg28a1c60t3gf');
      expect(tid.getType()).toBe('prefix');
    });

    it('should throw an error for invalid TypeID string', () => {
      const invalidStr = 'invalid_string_with_underscore';

      expect(() => {
        TypeID.fromString(invalidStr);
      }).toThrowError(new Error(`Invalid TypeID string: ${invalidStr}`));
    });
  });

  describe('fromUUIDBytes', () => {
    it('should construct TypeID from a UUID bytes without prefix', () => {
      const bytes = new Uint8Array([0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]);
      const tid = TypeID.fromUUIDBytes('', bytes);

      expect(tid.getSuffix()).toBe('00041061050r3gg28a1c60t3gf');
      expect(tid.getType()).toBe('');
    });

    it('should construct TypeID from a UUID bytes with prefix', () => {
      const bytes = new Uint8Array([0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]);
      const tid = TypeID.fromUUIDBytes('prefix', bytes);

      expect(tid.getSuffix()).toBe('00041061050r3gg28a1c60t3gf');
      expect(tid.getType()).toBe('prefix');
    });
  });

  describe('fromUUID', () => {
    it('should construct TypeID from a UUID string without prefix', () => {
      const uuid = "01889c89-df6b-7f1c-a388-91396ec314bc";
      const tid = TypeID.fromUUID('', uuid);

      expect(tid.getSuffix()).toBe('01h2e8kqvbfwea724h75qc655w');
      expect(tid.getType()).toBe('');
    });

    it('should construct TypeID from a UUID string  with prefix', () => {
      const uuid = "01889c89-df6b-7f1c-a388-91396ec314bc";
      const tid = TypeID.fromUUID('prefix', uuid);

      expect(tid.getSuffix()).toBe('01h2e8kqvbfwea724h75qc655w');
      expect(tid.getType()).toBe('prefix');
    });
  });

  describe('spec', () => {
    validJson.forEach((testcase: { name: string, prefix: string, typeid: string, uuid: string }) => {
      it(`should parse spec string: ${testcase.name}`, () => {
        const tid = TypeID.fromString(testcase.typeid);
        expect(tid.getType()).toBe(testcase.prefix);
        expect(tid.toString()).toBe(testcase.typeid);
        expect(tid.asUUID()).toBe(testcase.uuid);
      });

      it(`should encode spec uuid: ${testcase.name}`, () => {
        const tid = TypeID.fromUUID(testcase.prefix, testcase.uuid);
        expect(tid.getType()).toBe(testcase.prefix);
        expect(tid.toString()).toBe(testcase.typeid);
        expect(tid.asUUID()).toBe(testcase.uuid);
      });
    });

    invalidJson.forEach((testcase: { name: string, typeid: string }) => {
      it(`should throw an error for invalid spec test: ${testcase.name}`, () => {
        expect(() => {
          TypeID.fromString(testcase.typeid);
        }).toThrowError();
      });
    });
  });

});