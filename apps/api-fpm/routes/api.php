<?php

use App\Http\Controllers\BasketController;
use App\Http\Controllers\CatalogController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider and all of them will
| be assigned to the "api" middleware group. Make something great!
|
*/

Route::post('catalog/search', [CatalogController::class, 'search'])->name('catalog.search');
Route::get('basket/current', [BasketController::class, 'currentBasket'])->name('basket.current');
Route::post('basket/current/set-item', [BasketController::class, 'setItem'])->name('basket.set-item');
