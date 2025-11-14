<?php

use App\Http\Controllers\BasketController;
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

Route::get('baskets/current', [BasketController::class, 'getCurrentBasket'])->name('baskets.current');
Route::post('baskets/current/set-item', [BasketController::class, 'setItem'])->name('baskets.set-item');
