package config

import "time"

const SessionLifeTime = 24 * time.Hour
const AuthorityValidationTimeout = 5 * time.Minute
const AuthenticationTimeout = 3 * time.Minute
const CookieLifeTime = 7 * 24 * time.Hour
