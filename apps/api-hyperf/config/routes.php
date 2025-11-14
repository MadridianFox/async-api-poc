<?php

declare(strict_types=1);
/**
 * This file is part of Hyperf.
 *
 * @link     https://www.hyperf.io
 * @document https://hyperf.wiki
 * @contact  group@hyperf.io
 * @license  https://github.com/hyperf/hyperf/blob/master/LICENSE
 */

use App\Controller\BasketController;
use App\Controller\CatalogController;
use Hyperf\HttpServer\Router\Router;

Router::addRoute(['GET', 'POST', 'HEAD'], '/', 'App\Controller\IndexController@index');

Router::get('/favicon.ico', function () {
    return '';
});
Router::post('/api/catalog/search', [CatalogController::class, 'search']);
Router::get('/api/basket/current', [BasketController::class, 'currentBasket']);
Router::post('/api/basket/current/set-item', [BasketController::class, 'setItem']);