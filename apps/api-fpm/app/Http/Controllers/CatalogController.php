<?php

namespace App\Http\Controllers;

use App\Domains\Catalog\ProductClient;
use App\Http\Requests\SearchProductsRequest;
use App\Http\Resources\ProductSearchResultResource;
use Illuminate\Http\Resources\Json\JsonResource;

class CatalogController extends Controller
{
    public function search(SearchProductsRequest $request, ProductClient $client): JsonResource
    {
        return ProductSearchResultResource::make($client->searchProducts($request));
    }
}
