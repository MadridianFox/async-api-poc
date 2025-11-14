<?php

namespace App\Controller;

use App\Domain\Basket\BasketClient;
use App\Requests\CurrentBasketRequest;
use App\Requests\SetBasketItemRequest;
use App\Resources\BasketResource;
use Hyperf\HttpMessage\Exception\HttpException;

class BasketController extends AbstractController
{
    public function __construct(private readonly BasketClient $basketClient)
    {
        parent::__construct();
    }

    public function setItem(SetBasketItemRequest $request)
    {
        $basket = $this->basketClient->setItem(
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

    public function currentBasket(CurrentBasketRequest $request)
    {
        $basket = $this->basketClient->getCurrentBasket($request->getUserId(), $request->isWithItems());
        if (!$basket) {
            throw new HttpException(500, 'Internal Error');
        }

        return BasketResource::make($basket);
    }
}
