const fs = require('fs');
const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));

const endpoints = [];

for (const path in swagger.paths) {
    for (const method in swagger.paths[path]) {
        const details = swagger.paths[path][method];
        endpoints.push({
            method: method.toUpperCase(),
            path: path,
            summary: details.summary || '',
            description: details.description || '',
            tags: details.tags || []
        });
    }
}

console.log(JSON.stringify(endpoints, null, 2));
