const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));
const frontendDir = path.resolve('frontend');

const results = {
    timestamp: new Date().toISOString(),
    summary: {
        totalEndpoints: 0,
        integratedEndpoints: 0,
        missingEndpoints: 0,
        skippedEndpoints: 0
    },
    endpoints: [],
    fieldUsage: {}
};

// 從 swagger 提取所有 schema 定義
function extractAllSchemas() {
    const schemas = {};

    if (swagger.definitions) {
        for (const [name, def] of Object.entries(swagger.definitions)) {
            schemas[name] = extractFieldsFromSchema(def, name);
        }
    }

    return schemas;
}

// 遞迴提取 schema 中的欄位
function extractFieldsFromSchema(def, name, prefix = '') {
    const fields = {};

    if (!def) return fields;

    if (def.$ref) {
        const refName = def.$ref.split('/').pop();
        return { _ref: refName };
    }

    if (def.type === 'object' && def.properties) {
        for (const [fieldName, fieldDef] of Object.entries(def.properties)) {
            const fullName = prefix ? `${prefix}.${fieldName}` : fieldName;

            if (fieldDef.$ref) {
                const refName = fieldDef.$ref.split('/').pop();
                fields[fieldName] = { _ref: refName, _path: fullName };
            } else if (fieldDef.type === 'object') {
                fields[fieldName] = extractFieldsFromSchema(fieldDef, name, fullName);
            } else if (fieldDef.type === 'array' && fieldDef.items) {
                if (fieldDef.items.$ref) {
                    const refName = fieldDef.items.$ref.split('/').pop();
                    fields[fieldName] = { _arrayOf: refName, _path: fullName };
                } else {
                    fields[fieldName] = { _array: extractFieldsFromSchema(fieldDef.items, name, fullName), _path: fullName };
                }
            } else {
                fields[fieldName] = { type: fieldDef.type, _path: fullName, required: def.required?.includes(fieldName) };
            }
        }
    }

    return fields;
}

// 在 .vue 和 .ts 檔案中搜尋屬性存取
function findPropertyAccess() {
    const accessPatterns = {};

    // 搜尋 .vue 檔案中的屬性存取
    try {
        const vueFiles = execSync(`rg -g "*.vue" -n "\\.[a-zA-Z_][a-zA-Z0-9_]*" "${frontendDir}"`, {
            encoding: 'utf8',
            stdio: ['pipe', 'pipe', 'pipe']
        });

        const vueLines = vueFiles.split('\n');
        for (const line of vueLines) {
            if (!line.trim()) continue;

            const match = line.match(/^(.+):(\d+):(.+)/);
            if (match) {
                const [, filePath, lineNum, content] = match;
                const fileName = path.basename(filePath);

                const propertyMatches = content.matchAll(/\.([a-zA-Z_][a-zA-Z0-9_]*)/g);
                for (const pm of propertyMatches) {
                    const prop = pm[1];
                    const key = `${fileName}:${prop}`;

                    if (!accessPatterns[key]) {
                        accessPatterns[key] = {
                            file: fileName,
                            property: prop,
                            occurrences: [],
                            lines: []
                        };
                    }

                    accessPatterns[key].occurrences.push({
                        file: filePath,
                        line: parseInt(lineNum),
                        context: content.trim().substring(0, 100)
                    });
                    accessPatterns[key].lines.push(parseInt(lineNum));
                }
            }
        }
    } catch (e) {
        console.log('警告: 搜尋 Vue 檔案屬性存取時發生錯誤');
    }

    // 搜尋 .ts 檔案中的屬性存取
    try {
        const tsFiles = execSync(`rg -g "*.ts" -g "*.js" -n "\\.[a-zA-Z_][a-zA-Z0-9_]*" "${frontendDir}/composables" "${frontendDir}/stores" "${frontendDir}/types"`, {
            encoding: 'utf8',
            stdio: ['pipe', 'pipe', 'pipe']
        });

        const tsLines = tsFiles.split('\n');
        for (const line of tsLines) {
            if (!line.trim()) continue;

            const match = line.match(/^(.+):(\d+):(.+)/);
            if (match) {
                const [, filePath, lineNum, content] = match;
                const fileName = path.basename(filePath);

                const propertyMatches = content.matchAll(/\.([a-zA-Z_][a-zA-Z0-9_]*)/g);
                for (const pm of propertyMatches) {
                    const prop = pm[1];
                    const key = `${fileName}:${prop}`;

                    if (!accessPatterns[key]) {
                        accessPatterns[key] = {
                            file: fileName,
                            property: prop,
                            occurrences: [],
                            lines: []
                        };
                    }

                    accessPatterns[key].occurrences.push({
                        file: filePath,
                        line: parseInt(lineNum),
                        context: content.trim().substring(0, 100)
                    });
                    accessPatterns[key].lines.push(parseInt(lineNum));
                }
            }
        }
    } catch (e) {
        console.log('警告: 搜尋 TypeScript 檔案屬性存取時發生錯誤');
    }

    return accessPatterns;
}

// 搜尋常見的資料存取模式
function findDataAccessPatterns() {
    const patterns = {};

    const patternsToSearch = [
        { regex: /item\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'item' },
        { regex: /data\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'data' },
        { regex: /res\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'res' },
        { regex: /response\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'response' },
        { regex: /schedule\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'schedule' },
        { regex: /rule\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'rule' },
        { regex: /exception\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'exception' },
        { regex: /teacher\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'teacher' },
        { regex: /center\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'center' },
        { regex: /offering\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'offering' },
        { regex: /course\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'course' },
        { regex: /room\.([a-zA-Z_][a-zA-Z0-9_]*)/g, prefix: 'room' },
    ];

    for (const pattern of patternsToSearch) {
        try {
            const output = execSync(`rg -g "*.vue" -g "*.ts" -g "*.js" -n "${pattern.regex.source.replace(/\\/g, '\\\\')}" "${frontendDir}"`, {
                encoding: 'utf8',
                stdio: ['pipe', 'pipe', 'pipe']
            });

            const lines = output.split('\n');
            for (const line of lines) {
                if (!line.trim()) continue;

                const match = line.match(/^(.+):(\d+):(.+)/);
                if (match) {
                    const [, filePath, lineNum, content] = match;
                    const fileName = path.basename(filePath);

                    const matches = content.matchAll(pattern.regex);
                    for (const m of matches) {
                        const prop = m[1];
                        const key = `${pattern.prefix}:${prop}`;

                        if (!patterns[key]) {
                            patterns[key] = {
                                prefix: pattern.prefix,
                                property: prop,
                                files: new Set(),
                                occurrences: []
                            };
                        }

                        patterns[key].files.add(fileName);
                        patterns[key].occurrences.push({
                            file: filePath,
                            line: parseInt(lineNum),
                            context: content.trim().substring(0, 100)
                        });
                    }
                }
            }
        } catch (e) {
            // 忽略找不到的情況
        }
    }

    for (const key of Object.keys(patterns)) {
        patterns[key].files = Array.from(patterns[key].files);
    }

    return patterns;
}

// 搜尋 API 響應處理模式
function findApiResponseUsage() {
    const apiUsage = {
        responseFields: {},
        codeAccess: {},
        destructuring: []
    };

    // 搜尋 api().get/post 等調用
    try {
        const apiCalls = execSync(`rg -g "*.vue" -g "*.ts" -n "api\\([^)]+\\)" "${frontendDir}"`, {
            encoding: 'utf8',
            stdio: ['pipe', 'pipe', 'pipe']
        });

        const lines = apiCalls.split('\n');
        for (const line of lines) {
            if (!line.trim()) continue;

            const match = line.match(/^(.+):(\d+):(.+)/);
            if (match) {
                const [, filePath, lineNum, content] = match;
                apiUsage.apiCalls = apiUsage.apiCalls || [];
                apiUsage.apiCalls.push({
                    file: path.basename(filePath),
                    line: parseInt(lineNum),
                    code: content.trim().substring(0, 150)
                });
            }
        }
    } catch (e) {
        // 忽略找不到的情況
    }

    // 搜尋 .data 或 .datas 存取（API 響應中的資料欄位）
    try {
        const dataAccess = execSync(`rg -g "*.vue" -g "*.ts" -n "\\.data[s]?\\." "${frontendDir}"`, {
            encoding: 'utf8',
            stdio: ['pipe', 'pipe', 'pipe']
        });

        const lines = dataAccess.split('\n');
        for (const line of lines) {
            if (!line.trim()) continue;

            const match = line.match(/^(.+):(\d+):(.+)/);
            if (match) {
                const [, filePath, lineNum, content] = match;
                const subMatches = content.matchAll(/\.data[s]?\.([a-zA-Z_][a-zA-Z0-9_]*)/g);

                for (const sm of subMatches) {
                    const field = sm[1];
                    apiUsage.responseFields[field] = apiUsage.responseFields[field] || {
                        count: 0,
                        files: new Set(),
                        lines: []
                    };
                    apiUsage.responseFields[field].count++;
                    apiUsage.responseFields[field].files.add(path.basename(filePath));
                    apiUsage.responseFields[field].lines.push(parseInt(lineNum));
                }
            }
        }
    } catch (e) {
        // 忽略找不到的情況
    }

    // 轉換 Set 為 Array
    for (const field of Object.keys(apiUsage.responseFields)) {
        apiUsage.responseFields[field].files = Array.from(apiUsage.responseFields[field].files);
    }

    return apiUsage;
}

for (const epPath in swagger.paths) {
    for (const method in swagger.paths[epPath]) {
        const details = swagger.paths[epPath][method];

        const searchPath = epPath.replace(/\{[^}]+\}/g, '').replace(/:[^/]+/g, '').replace(/\/$/, '');

        if (searchPath === '' || searchPath === '/api/v1') {
            results.summary.skippedEndpoints++;
            results.endpoints.push({
                path: epPath,
                method: method.toUpperCase(),
                integrated: 'SKIPPED',
                files: [],
                fields: []
            });
            continue;
        }

        results.summary.totalEndpoints++;

        try {
            let searchStr = searchPath;
            if (searchStr.startsWith('/api/v1')) {
                searchStr = searchStr.substring(7);
            }

            const output = execSync(`rg -l "${searchStr}" "${frontendDir}"`, { encoding: 'utf8' });
            const files = output.trim().split('\n').filter(f => f !== '');

            const usedFields = new Set();

            for (const file of files) {
                if (!file.match(/\.(vue|ts|js)$/)) continue;

                try {
                    const content = fs.readFileSync(file, 'utf8');

                    const propertyMatches = content.matchAll(/\.(id|name|type|status|date|time|at|code|message|data|datas)/gi);
                    for (const m of propertyMatches) {
                        usedFields.add(m[1].toLowerCase());
                    }

                    const modelProps = content.matchAll(/\.(offering_name|teacher_name|room_name|start_time|end_time|weekday|course_name|effective_range)/gi);
                    for (const m of modelProps) {
                        usedFields.add(m[1].toLowerCase());
                    }
                } catch (e) {
                    // 忽略讀取錯誤
                }
            }

            results.summary.integratedEndpoints++;
            results.endpoints.push({
                path: epPath,
                method: method.toUpperCase(),
                summary: details.summary,
                integrated: files.length > 0,
                files: files,
                fields: Array.from(usedFields)
            });
        } catch (e) {
            results.summary.missingEndpoints++;
            results.endpoints.push({
                path: epPath,
                method: method.toUpperCase(),
                summary: details.summary,
                integrated: false,
                files: [],
                fields: []
            });
        }
    }
}

results.fieldUsage = {
    allAccessPatterns: findPropertyAccess(),
    dataPatterns: findDataAccessPatterns(),
    apiResponseUsage: findApiResponseUsage()
};

const fieldStats = {};
for (const [key, data] of Object.entries(results.fieldUsage.dataPatterns)) {
    const [prefix, prop] = key.split(':');
    if (!fieldStats[prefix]) {
        fieldStats[prefix] = {
            properties: new Set(),
            totalOccurrences: 0
        };
    }
    fieldStats[prefix].properties.add(prop);
    fieldStats[prefix].totalOccurrences += data.occurrences.length;
}

results.fieldUsage.stats = {};
for (const [prefix, data] of Object.entries(fieldStats)) {
    results.fieldUsage.stats[prefix] = {
        uniqueProperties: data.properties.size,
        totalOccurrences: data.totalOccurrences,
        properties: Array.from(data.properties)
    };
}

// 提取 swagger schemas 以便比對
results.swaggerSchemas = extractAllSchemas();

fs.writeFileSync('integration_audit.json', JSON.stringify(results, null, 2), 'utf8');
console.log('Audit complete. Results saved to integration_audit.json');
console.log(`\nSummary:`);
console.log(`  Total endpoints: ${results.summary.totalEndpoints}`);
console.log(`  Integrated: ${results.summary.integratedEndpoints}`);
console.log(`  Missing: ${results.summary.missingEndpoints}`);
console.log(`  Skipped: ${results.summary.skippedEndpoints}`);
console.log(`\nField usage patterns detected: ${Object.keys(results.fieldUsage.dataPatterns).length}`);
console.log(`\nSwagger schemas extracted: ${Object.keys(results.swaggerSchemas).length}`);