<?php

namespace App\Http\Controllers;

use App\Http\Requests\SearchProductsRequest;
use App\Http\Resources\ProductResource;
use App\Models\Product;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Http\Resources\Json\JsonResource;

class ProductController extends Controller
{
    public function index(SearchProductsRequest $request): JsonResource
    {
        /** @var Builder $query */
        $query = Product::query();

        $minQty = $request->getQtyMin();
        if ($minQty !== null) {
            $query->where('qty', '>=', $minQty);
        }

        $maxQty = $request->getQtyMax();
        if ($maxQty !== null) {
            $query->where('qty', '<=', $maxQty);
        }

        $minPrice = $request->getPriceMin();
        if ($minPrice !== null) {
            $query->where('price', '>=', $minPrice);
        }

        $maxPrice = $request->getPriceMax();
        if ($maxPrice !== null) {
            $query->where('price', '<=', $maxPrice);
        }

        $name = $request->getNameLike();
        if ($name !== null) {
            $query->where('name', 'like', "%$name%");
        }

        return ProductResource::collection($query->paginate(20));
    }

    public function show(Product $product): JsonResource
    {
        return ProductResource::make($product);
    }
}
