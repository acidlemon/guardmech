require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  "root": true,
  "extends": [
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting',
  ],
  "parser": "vue-eslint-parser",
  "parserOptions": {
    "ecmaVersion": 'latest',
    "parser": "@typescript-eslint/parser",
    "sourceType": "module"
  },
  "env": {
    "es2021": true,
    'vue/setup-compiler-macros': true,
  },
  "globals": {},
  "rules": {
    "no-unused-vars": "off",
    "vue/no-multiple-template-root": "off",
    "@typescript-eslint/no-unused-vars" : [ "error", { "argsIgnorePattern": "^_", "varsIgnorePattern": "^_" } ]
  }
}
