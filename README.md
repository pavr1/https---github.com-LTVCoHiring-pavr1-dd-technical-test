# https---github.com-LTVCoHiring-pavr1-dd-technical-test

This is a small challege project created as part of a Golang hiring process which implies calling apis for songs research. 

Arguments:
-FROM: Initial date (YYYY-MM-DD)
-UNTIL: FInal date (YYYY-MM-DD)
-Artist: Artist Name (Optional)

At start up the app caches data from a specific date up to 90 days after, so the retrieval will not be slow, in case some specific date does not exist in the cache then 
it will be retireved from the API and saved in the cache for future retrieves.
