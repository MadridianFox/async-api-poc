<?php

namespace App\Domains\Catalog;

interface ProductFilter
{
    public function getMinQty(): ?int;
    public function getMaxQty(): ?int;
    public function getMinPrice(): ?int;
    public function getMaxPrice(): ?int;
    public function getNameLike(): ?string;
    public function getPage(): ?int;
}
