/**
 * Teacher Mock JWT Token Generator for Testing
 * 
 * This script generates a valid Teacher JWT token for testing Teacher APIs.
 * 
 * Usage: node generate-teacher-token.js
 */

const crypto = require('crypto');

function base64UrlEncode(str) {
  return Buffer.from(str)
    .toString('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
}

function sign(data, secret) {
  return crypto
    .createHmac('sha256', secret)
    .update(data)
    .digest('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
}

const secretKey = 'mock-secret-key-for-development';

// Teacher Token (User ID = 1, UserType = TEACHER)
const teacherHeader = {
  alg: 'HS256',
  typ: 'JWT'
};

const teacherClaims = {
  user_type: 'TEACHER',
  user_id: 1,  // teacher1@example.com from database
  center_id: 1,  // Will be populated from teacher's memberships
  exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60)  // 24 hours
};

const headerEncoded = base64UrlEncode(JSON.stringify(teacherHeader));
const claimsEncoded = base64UrlEncode(JSON.stringify(teacherClaims));
const teacherSignature = `${headerEncoded}.${claimsEncoded}`;

const teacherToken = `${headerEncoded}.${claimsEncoded}.${sign(`${headerEncoded}.${claimsEncoded}`, secretKey)}`;

console.log('========================================');
console.log('Teacher Mock JWT Token');
console.log('========================================');
console.log();
console.log('Teacher ID: 1 (teacher1@example.com)');
console.log('UserType: TEACHER');
console.log();
console.log('Token:');
console.log(teacherToken);
console.log();

// Admin Token for comparison
const adminClaims = {
  user_type: 'ADMIN',
  user_id: 16,  // admin@timeledger.com
  center_id: 1,
  exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60)
};

const adminClaimsEncoded = base64UrlEncode(JSON.stringify(adminClaims));
const adminToken = `${headerEncoded}.${adminClaimsEncoded}.${sign(`${headerEncoded}.${adminClaimsEncoded}`, secretKey)}`;

console.log('========================================');
console.log('Admin Mock JWT Token (for comparison)')
console.log('========================================');
console.log();
console.log('Admin ID: 16 (admin@timeledger.com)');
console.log('UserType: ADMIN');
console.log();
console.log('Token:');
console.log(adminToken);
console.log();

// Export for use in other scripts
module.exports = { teacherToken, adminToken };
