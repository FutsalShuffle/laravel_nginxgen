package common

// Phpdoc params
const EnumParamRegex = `@NGEnum (.+?( ))(.+)`
const IntParamRegex = `@NGIntOnly (.+?( |\n))`
const StringParamRegex = `@NGStringOnly (.+?( |\n))`
const LimitQueryParamRegex = `@NGQLimit (.+?( |\n))`

// nginx
const (
	NginxGateHandle     = "try_files $uri $uri/ /index.php?$query_string;"
	RandomLetterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NginxIntOnlyPath    = "([0-9]*)"
	NginxStringOnlyPath = "(\\D+)"
	PDocParamSeparator  = ","
)

// laravel
const LaravelCacheRoutePath = "/bootstrap/cache/"
