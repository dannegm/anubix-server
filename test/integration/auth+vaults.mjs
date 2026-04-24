import crypto from 'crypto';

const BASE_URL = `http://localhost:${process.env.PORT ?? 8080}`;

console.log('Running integration tests against:', BASE_URL);

// ── Helpers ──────────────────────────────────────────────────────────────────

function generateSalt() {
    return crypto.randomBytes(16).toString('hex');
}

function deriveAuthHash(password, salt) {
    return crypto.pbkdf2Sync(password, salt, 100000, 32, 'sha256').toString('hex');
}

function generateVaultKey() {
    return crypto.randomBytes(32);
}

function encryptVaultKey(vaultKey, masterKey) {
    const iv = crypto.randomBytes(12);
    const cipher = crypto.createCipheriv(
        'aes-256-gcm',
        Buffer.from(masterKey, 'hex').slice(0, 32),
        iv,
    );
    const encrypted = Buffer.concat([cipher.update(vaultKey), cipher.final()]);
    const authTag = cipher.getAuthTag();
    return {
        encryptedVaultKey: encrypted.toString('base64'),
        iv: iv.toString('base64'),
        authTag: authTag.toString('base64'),
    };
}

async function request(method, path, body, token) {
    const res = await fetch(`${BASE_URL}${path}`, {
        method,
        headers: {
            'Content-Type': 'application/json',
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: body ? JSON.stringify(body) : undefined,
    });

    if (res.status === 204) return null; // No content
    const data = await res.json();
    if (!res.ok) throw new Error(`${method} ${path} → ${res.status}: ${JSON.stringify(data)}`);
    return data;
}

// ── Flow ─────────────────────────────────────────────────────────────────────

const email = `test_${Date.now()}@anubix.dev`;
const password = 'supersecretpassword';

console.log('─'.repeat(50));
console.log('🔐 ANUBIX AUTH FLOW TEST');
console.log('─'.repeat(50));

// 1. Register
console.log('\n[1] Registering user...');
const salt = generateSalt();
const masterKey = deriveAuthHash(password, salt);
const authHash = deriveAuthHash(masterKey, salt + '_auth');
const vaultKey = generateVaultKey();
const { encryptedVaultKey, iv, authTag } = encryptVaultKey(vaultKey, masterKey);

const registerRes = await request('POST', '/auth/register', {
    email,
    auth_hash: authHash,
    salt,
    vault_name: 'Personal',
    encrypted_vault_key: encryptedVaultKey,
    vault_key_iv: iv,
    vault_key_auth_tag: authTag,
});
console.log('✅ Registered:', registerRes);

// 2. Get salt
console.log('\n[2] Fetching salt...');
const saltRes = await request('GET', `/auth/salt?email=${email}`);
console.log('✅ Salt:', saltRes);

// 3. Login
console.log('\n[3] Logging in...');
const loginMasterKey = deriveAuthHash(password, saltRes.salt);
const loginAuthHash = deriveAuthHash(loginMasterKey, saltRes.salt + '_auth');
const loginRes = await request('POST', '/auth/login', {
    email,
    auth_hash: loginAuthHash,
});
console.log('✅ Login token (no device):', loginRes.token);
const tokenNoDevice = loginRes.token;

// 4. Register device
console.log('\n[4] Registering device...');
const deviceRes = await request(
    'POST',
    '/devices',
    {
        name: 'MacBook Pro',
        fingerprint: crypto.randomBytes(16).toString('hex'),
        device_type: 'web',
    },
    tokenNoDevice,
);
console.log('✅ Device:', deviceRes);

// 5. Get token with device
console.log('\n[5] Getting token with device_id...');
const tokenRes = await request('POST', '/auth/token', { device_id: deviceRes.id }, tokenNoDevice);
console.log('✅ Token with device:', tokenRes.token);
const token = tokenRes.token;

// 6. /me
console.log('\n[6] GET /me...');
const me = await request('GET', '/me', null, token);
console.log('✅ Me:', JSON.stringify(me, null, 2));

// 7. GET /vaults
console.log('\n[7] GET /vaults...');
const vaults = await request('GET', '/vaults', null, token);
console.log('✅ Vaults:', JSON.stringify(vaults, null, 2));

// 8. POST /vaults
console.log('\n[8] Creating new vault...');
const newVaultKey = generateVaultKey();
const newEncrypted = encryptVaultKey(newVaultKey, masterKey);
const newVault = await request(
    'POST',
    '/vaults',
    {
        name: 'Work',
        encrypted_vault_key: newEncrypted.encryptedVaultKey,
        vault_key_iv: newEncrypted.iv,
        vault_key_auth_tag: newEncrypted.authTag,
    },
    token,
);
console.log('✅ New vault:', JSON.stringify(newVault, null, 2));

// 9. GET /vaults/:id
console.log('\n[9] GET /vaults/:id...');
const vault = await request('GET', `/vaults/${newVault.id}`, null, token);
console.log('✅ Vault by ID:', JSON.stringify(vault, null, 2));

// 10. PUT /vaults/:id
console.log('\n[10] PUT /vaults/:id...');
const updatedVault = await request(
    'PUT',
    `/vaults/${newVault.id}`,
    {
        name: 'Work Updated',
        encrypted_vault_key: newEncrypted.encryptedVaultKey,
        vault_key_iv: newEncrypted.iv,
        vault_key_auth_tag: newEncrypted.authTag,
    },
    token,
);
console.log('✅ Updated vault:', JSON.stringify(updatedVault, null, 2));

// 11. DELETE /vaults/:id
console.log('\n[11] DELETE /vaults/:id...');
await request('DELETE', `/vaults/${newVault.id}`, null, token);
console.log('✅ Vault deleted');

console.log('\n' + '─'.repeat(50));
console.log('✅ ALL TESTS PASSED');
console.log('─'.repeat(50));
