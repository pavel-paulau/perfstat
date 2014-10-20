perfstat
========

[![GoDoc](https://godoc.org/github.com/pavel-paulau/perfstat?status.svg)](https://godoc.org/github.com/pavel-paulau/perfstat)

**perfstat** is a plugable collector of miscellaneous stats. It works similar to [dstat](https://github.com/dagwieers/dstat) but distributed as a binary executable and also introduces more flexible plugin development.

Most importantly it allows to asynchronously store samples in [perfkeeper](https://github.com/pavel-paulau/perfkeeper)

Installation and usage
----------------------

Download the latest executable file at [Releases](https://github.com/pavel-paulau/perfstat/releases) page.

Alternatively get [Go](http://golang.org/doc/install) and then get perfstat:

    $ go get github.com/pavel-paulau/perfstat

Now let's sample CPU usage:

	> perfstat -cpu
	cpu_user cpu_sys cpu_idle cpu_iowait 
	------------------------------------
	     127       9      264          0 
	     105       4      291          0 
	     105       3      292          0 
	     108       1      291          0 
	     108       4      288          0 

How about storing memory stats in perfkeeper?

	> perfstat -mem -snapshot="mybenchmark" -source="kernel"
	mem_used mem_free mem_buff mem_cache 
	------------------------------------
	    2848     2306      242      2053 
	    2858     2290      242      2059 
	    2857     2291      242      2059 
	    2858     2290      242      2059 

`-h` will help to understand the other options.
