memrecycle ♻
==========

Why waste memory when you can recycle it?

All of the initial code and concepts were taken from this article on CloudFlare's blog:
http://blog.cloudflare.com/recycling-memory-buffers-in-go

The basic idea is that garbage collected languages like Go tend to over-request memory and become wasteful in certain scenarios (great graphs are provided in the article linked above that visualize this fact). This method introduces a way of recycling memory that yields truly staggering results by using goroutines and channels. The best example of this is when running the garbage_creator.go program for 10 minutes there are roughly 600 buffers created, 375,000,000 bytes requested, and only 150,000,000 bytes truly needed. So the program requested roughly 2.5x the amount of memory it actually needed. Additionally, there is an incredible amount of memory being used that is just idling waiting to be taken care of by the garbage collector. When running the improved mem_manage.go program for 10 minutes there are roughly 21 buffers created, 160,000,000 bytes requested, and only 150,000,000 bytes truly needed. This is almost a 1:1 ratio between what is required and what is requested - a drastic improvement. There is also close to 0 idling memory.

