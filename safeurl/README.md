#safeUrl

	unescape url safely


##usage same as go standard url


##description

	an url encoded string like

	aabb%A 
	aabb%
	aabb%Accc
	aabb%AAcc
	aabb%2cc

	will return an error

	my solution:

	replace these undecoded '%' as '_'	


##TODO:

update support unescape path