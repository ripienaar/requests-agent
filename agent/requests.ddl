metadata :name        => "requests",
         :description => "Network requests agent",
         :author      => "R.I.Pienaar <rip@devco.net>",
         :license     => "Apache-2.0",
         :version     => "0.0.1",
         :url         => "http://choria.io",
         :provider    => "external",
         :timeout     => 120


action "request", :description => "Performs a HTTP request" do
  display :always

  input :body,
        :prompt      => "Body",
        :description => "The body to send in the request",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 4096,
        :optional    => true


  input :body_file,
        :prompt      => "Body File",
        :description => "Sends a specific file locally to each node as body",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 512,
        :optional    => true


  input :headers,
        :prompt      => "Headers",
        :description => "Request headers to send as k=v pairs",
        :type        => :hash,
        :optional    => true


  input :method,
        :prompt      => "Method",
        :description => "The HTTP method to use",
        :type        => :list,
        :default     => "GET",
        :list        => ["GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS", "HEAD"],
        :optional    => true


  input :password,
        :prompt      => "Password",
        :description => "Password for authentication",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 128,
        :optional    => true


  input :query,
        :prompt      => "Query",
        :description => "Request query to send as k=v pairs",
        :type        => :hash,
        :optional    => true


  input :statuscode,
        :prompt      => "Expected Statuscode",
        :description => "Checks the result against this statuscode, else expects 200",
        :type        => :integer,
        :optional    => true


  input :url,
        :prompt      => "URL",
        :description => "Address to fetch",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 512,
        :optional    => false


  input :username,
        :prompt      => "Username",
        :description => "Username for authentication",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 128,
        :optional    => true




  output :body,
         :description => "The response body",
         :type        => "string",
         :display_as  => "Body"

  output :duration,
         :description => "How long, in seconds, the request took to complete",
         :type        => "float",
         :display_as  => "Duration"

  output :headers,
         :description => "The response headers",
         :type        => "hash",
         :display_as  => "Headers"

  output :statuscode,
         :description => "The response statuscode",
         :type        => "integer",
         :display_as  => "Statuscode"

  summarize do
    aggregate summary(:statuscode)
  end
end

action "download", :description => "Downloads a file" do
  display :failed

  input :headers,
        :prompt      => "Headers",
        :description => "Request headers to send as k=v pairs",
        :type        => :hash,
        :optional    => true


  input :md5,
        :prompt      => "Expected MD5",
        :description => "If set only files matching this digest will be saved",
        :type        => :string,
        :validation  => '^[0-9a-z]+$',
        :maxlength   => 32,
        :optional    => true


  input :password,
        :prompt      => "Password",
        :description => "Password for authentication",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 128,
        :optional    => true


  input :query,
        :prompt      => "Query",
        :description => "Request query to send as k=v pairs",
        :type        => :hash,
        :optional    => true


  input :target,
        :prompt      => "Target File",
        :description => "File to save the response into",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 512,
        :optional    => false


  input :target_mode,
        :prompt      => "File Mode",
        :description => "Mode to assign to the file post download in octal numeric format",
        :type        => :string,
        :validation  => '^\d+$',
        :maxlength   => 4,
        :optional    => true


  input :url,
        :prompt      => "URL",
        :description => "Address to fetch",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 512,
        :optional    => false


  input :username,
        :prompt      => "Username",
        :description => "Username for authentication",
        :type        => :string,
        :validation  => '.+',
        :maxlength   => 128,
        :optional    => true




  output :bytes,
         :description => "Number of bytes fetched",
         :type        => "integer",
         :display_as  => "Bytes"

  output :error,
         :description => "Error string on failure",
         :type        => "string",
         :display_as  => "Error"

  output :md5,
         :description => "MD5 digest of the downloaded file",
         :type        => "string",
         :display_as  => "MD5"

  output :statuscode,
         :description => "HTTP Status Code",
         :type        => "integer",
         :display_as  => "Status"

  summarize do
    aggregate summary(:bytes)
    aggregate summary(:md5)
    aggregate summary(:statuscode)
  end
end

