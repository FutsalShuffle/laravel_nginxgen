Параметры:
 - @NGEnum {параметр из роута без скобок} Значения,через,запятую
 - @NGIntOnly {параметр из роута без скобок}
 - @NGStringOnly {параметр из роута без скобок}
 - @NGRLimit разрешенные,параметры,query \
Назначаются в методе контроллера **внутри phpdoc** \
Перед генерацией нужно сперва запустить **artisan route:cache**, без существующего кеша отдаст стандартную конфигурацию nginx для laravel \
Пример phpdoc для метода:

```php
 /**
     * @NGIntOnly id
     * @throws Exception
     * @param int $id
     * @param Request $request
     */
    public function downloadZip(int $id, Request $request): \Symfony\Component\HttpFoundation\BinaryFileResponse|RedirectResponse
    {
```
web.api
```php
    Route::get('download-zip/{id}',
        [ZipController::class, 'downloadZip']);
```
Сгенерирует такое правило nginx 
```
location ~ /download-zip/([0-9]*) {
   try_files $uri $uri/ /index.php?$query_string;
}
```
Использование билда: 
```
./laravel_nginxgen -project=path -output=path.conf

По дефолту можно закинуть в корень проекта (пути ., ./locations.conf)
```