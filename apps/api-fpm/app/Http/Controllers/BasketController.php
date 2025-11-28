<?php

namespace App\Http\Controllers;

use App\Domains\Basket\BasketClient;
use App\Http\Requests\CurrentBasketRequest;
use App\Http\Requests\SetBasketItemRequest;
use App\Http\Resources\BasketResource;
use Symfony\Component\HttpKernel\Exception\HttpException;

class BasketController extends Controller
{
    public function setItem(SetBasketItemRequest $request, BasketClient $basketClient)
    {
        $basket = $basketClient->setItem(
            $request->getUserId(),
            $request->getProductId(),
            $request->getQty(),
            $request->isWithItems()
        );

        if (!$basket) {
            throw new HttpException(500, 'Internal Error');
        }

        return BasketResource::make($basket);
    }

    public function currentBasket(CurrentBasketRequest $request, BasketClient $basketClient)
    {
        $basket = $basketClient->getCurrentBasket($request->getUserId(), $request->isWithItems());
        if (!$basket) {
            throw new HttpException(500, 'Internal Error');
        }

        return BasketResource::make($basket);
    }
}
