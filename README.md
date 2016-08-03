# Serial

## Serial number and nonce generation in Go

This is a fairly trivial library to generate guaranteed unique serial numbers
as int64 values, while ensuring thread safety.

It also provides functions to let you flag which values you've seen, so you
can use them as nonce values. An example use case would be as `jti` parameters
for JWT authorization tokens.

If you use the 'seen' flag feature, remember to expire the history
periodically, based on the lifetime of your tokens.

If you found yourself here because you want a library for JWT, JWS and JWE, I
suggest <https://github.com/lestrrat/go-jwx>.

Example usage:

    jti := serial.Generate()
    
    serial.SetSeen(jti)

    if serial.Seen(jti) {
      // Attempt to reuse nonce
    }

    // Tokens we issue have an exp value of half an hour from moment of issue,
    // so every now and again we...
    serial.ExpireHistory(time.Minute * 30)

