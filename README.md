Phpdoc parameters:
 - @NGEnum {route_param_without_{}} values,separated,by,comma
 - @NGIntOnly {route_param_without_{}}
 - @NGStringOnly {route_param_without_{}}
 - @NGQLimit values,separated,by,comma \
You can set these parameters for a controller method **inside phpdoc** \
You need to run **artisan route:cache** first, it will generate a default config unless you cache your routes \
Example of a phpdoc:

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
Generated Nginx rule
```
location ~ /download-zip/([0-9]*) {
   if ($request_method !~ ^(GET|HEAD)$) {
        return 404;
   }
   try_files $uri $uri/ /index.php?$query_string;
}
```
Build usage:
```
./laravel_nginxgen -project=path -output=path.conf

Parameters:
 -project (default: .)
 -output (default: ./locations.conf)
 -nginx-handle (default: try_files $uri $uri/ /index.php?$query_string;) - you will need to change this if you're using Laravel Octane
 -nginx-wms (default: 404) - Response status on a wrong method call
 -nginx-wqps (default: 404) - Response status on a wrong query parameters (when using @NGQLimit)
 -add-post (default: y) - Set N to disable POST addition for route methods when PUT/PATCH is present (by default web routes use POST instead with _method in body)
```
