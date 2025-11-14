<?php

namespace App\Controller;

use App\Domain\Catalog\CatalogClient;
use App\Requests\SearchProductsRequest;
use App\Resources\ProductSearchResultResource;
use Hyperf\Resource\Json\JsonResource;

class CatalogController extends AbstractController
{
    public function __construct(private readonly CatalogClient $client)
    {
        parent::__construct();
    }

    public function search(SearchProductsRequest $request): JsonResource
    {
        return ProductSearchResultResource::make($this->client->searchProducts($request));
    }
}
