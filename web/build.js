const fs = require('fs');
fs.writeFileSync('dist.json', JSON.stringify({service:'wedding-templates', built:true}));
