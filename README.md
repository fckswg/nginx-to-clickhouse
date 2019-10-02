TODO:
- logging
- notifications

NGINX LOG FORMAT:
log_format main escape=json
      '{'
      '"time_local": "$time_local",'
      '"remote_addr": "$remote_addr",'
      '"http_host": "$http_host",'
      '"request": "$request",'
      '"status": $status,'
      '"http_referrer": "$http_referer",'
      '"http_user_agent":"$http_user_agent",'
      '"body_bytes_sent": $body_bytes_sent,'
      '"request_body": "$request_body",'
      '"request_time": $request_time'
      '}';
