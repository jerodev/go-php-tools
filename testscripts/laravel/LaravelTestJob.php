<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;

class LaravelTestJob implements ShouldQueue
{
    use InteractsWithQueue, Queueable;

    public string $contents;

    public function handle(): void
    {
        \file_put_contents(__DIR__ . '/foo_bar', \implode(',' $this->contents));
    }
}
