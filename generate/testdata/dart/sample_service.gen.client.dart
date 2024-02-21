// Code generated by Frodo - DO NOT EDIT.
//
//   Timestamp: Wed, 21 Feb 2024 11:49:20 EST
//   Source:    sample_service.go
//   Generator: https://github.com/bridgekitio/frodo
//
import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;

class SampleServiceClient {
  final String baseURL;
  String authorization;
  http.Client httpClient = http.Client();

  SampleServiceClient(this.baseURL, {
      this.authorization = '',
  });

  
  /// Authorization regurgitates the "Authorization" metadata/header.
  Future<SampleResponse> Authorization(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Authorization';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// ComplexValues flexes our ability to encode/decode non-flat structs.
  Future<SampleComplexResponse> ComplexValues(SampleComplexRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.ComplexValues';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleComplexResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// ComplexValuesPath flexes our ability to encode/decode non-flat structs while
  /// specifying them via path and query string.
  Future<SampleComplexResponse> ComplexValuesPath(SampleComplexRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/complex/values/{InUser.ID}/{InUser.Name}/woot';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleComplexResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// CustomRoute performs a service operation where you override default behavior
  /// by providing routing-related Doc Options.
  Future<SampleResponse> CustomRoute(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/custom/route/1/{ID}/{Text}';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// CustomRouteBody performs a service operation where you override default behavior
  /// by providing routing-related Doc Options, but rely on body encoding rather than path.
  Future<SampleResponse> CustomRouteBody(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'PUT';
    var route = '/v2/custom/route/3/{ID}';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// CustomRouteQuery performs a service operation where you override default behavior
  /// by providing routing-related Doc Options. The input data relies on the path
  Future<SampleResponse> CustomRouteQuery(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/custom/route/2/{ID}';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Defaults simply utilizes all of the framework's default behaviors.
  Future<SampleResponse> Defaults(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Defaults';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Download results in a raw stream of data rather than relying on auto-encoding
  /// the response value.
  Future<SampleDownloadResponse> Download(SampleDownloadRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/download';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      var stream = await _handleResponseStream(response);
      return SampleDownloadResponse(
        content: stream.content,
        contentLength: stream.contentLength,
        contentType: stream.contentType,
        contentFileName: stream.contentFileName,
        contentRange: stream.contentRange,
      );
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// DownloadResumable results in a raw stream of data rather than relying on auto-encoding
  /// the response value. The stream includes Content-Range info as though you could resume
  /// your stream/download progress later.
  Future<SampleDownloadResponse> DownloadResumable(SampleDownloadRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/download/resumable';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      var stream = await _handleResponseStream(response);
      return SampleDownloadResponse(
        content: stream.content,
        contentLength: stream.contentLength,
        contentType: stream.contentType,
        contentFileName: stream.contentFileName,
        contentRange: stream.contentRange,
      );
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Fail4XX always returns a non-nil 400-series error.
  Future<SampleResponse> Fail4XX(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Fail4XX';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Fail5XX always returns a non-nil 500-series error.
  Future<SampleResponse> Fail5XX(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Fail5XX';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// ListenerA fires on only one of the triggers.
  Future<SampleResponse> ListenerA(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/ListenerA/Woot';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// ListenerB fires on multiple triggers... including another event-based endpoint. We also
  /// listen for the TriggerFailure event which should never fire properly.
  Future<SampleResponse> ListenerB(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    throw new SampleServiceException(501, 'ListenerB is not supported in the API gateway');
  }
  
  /// OmitMe exists in the service, but should be excluded from the public API.
  Future<SampleResponse> OmitMe(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    throw new SampleServiceException(501, 'OmitMe is not supported in the API gateway');
  }
  
  /// Panic um... panics. It never succeeds. It always behaves like me when I'm on a high place looking down.
  Future<SampleResponse> Panic(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Panic';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Redirect results in a 307-style redirect to the Download endpoint.
  Future<SampleRedirectResponse> Redirect(SampleRedirectRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/redirect';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      var stream = await _handleResponseStream(response);
      return SampleRedirectResponse(
        content: stream.content,
        contentLength: stream.contentLength,
        contentType: stream.contentType,
        contentFileName: stream.contentFileName,
        contentRange: stream.contentRange,
      );
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// SecureWithRoles lets us test role based security by looking at the 'roles' doc option.
  Future<SampleSecurityResponse> SecureWithRoles(SampleSecurityRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.SecureWithRoles';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleSecurityResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// Sleep successfully responds, but it will sleep for 5 seconds before doing so. Use this
  /// for test cases where you want to try out timeouts.
  Future<SampleResponse> Sleep(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.Sleep';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  Future<SampleResponse> TriggerFailure(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.TriggerFailure';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  Future<SampleResponse> TriggerLowerCase(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/v2/SampleService.TriggerLowerCase';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  
  /// TriggerUpperCase ensures that events still fire as "SampleService.TriggerUpperCase" even though
  /// we are going to set a different HTTP path.
  Future<SampleResponse> TriggerUpperCase(SampleRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'GET';
    var route = '/v2/Upper/Case/WootyAndTheBlowfish';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => SampleResponse.fromJson(json));
    } on SampleServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw SampleServiceException(500, e.toString());
    }
  }
  

  String _buildRequestPath(String method, String route, Map<String, dynamic> requestJson) {
    String stringify(Map<String, dynamic> json, String key) {
      return Uri.encodeComponent(json[key]?.toString() ?? '');
    }
    String stringifyAndRemove(Map<String, dynamic> json, String key) {
      return Uri.encodeComponent(json.remove(key)?.toString() ?? '');
    }

    // Since we're embedding values in a path or query string, we need to flatten "{a: {b: {c: 4}}}"
    // down to "a.b.c=4" for it to fit nicely into our URL-based binding.
    requestJson = _flattenJson(requestJson);

    // Replace variable segments w/ their values (e.g "/user/{User.ID}/file/{ID}" -> "/user/123/file/456").
    var resolvedPath = route
      .split('/')
      .map((s) => _isParameterSegment(s) ? stringifyAndRemove(requestJson, s.substring(1, s.length - 1)) : s)
      .join('/');

    // These encode the data in the body, so no need to shove it in the query string.
    if (method == 'POST' || method == 'PUT' || method == 'PATCH') {
      return resolvedPath;
    }

    // GET/DELETE/etc will pass all values through the query string.
    var queryValues = requestJson.keys
      .map((key) => key + '=' + stringify(requestJson, key))
      .join('&');

    return resolvedPath + '?' + queryValues;
  }

  bool _isParameterSegment(String segment) {
    return segment.startsWith('{') && segment.endsWith('}');
  }

  Future<T> _handleResponseJson<T>(http.StreamedResponse response, T Function(Map<String, dynamic>) factory) async {
    if (response.statusCode >= 400) {
      throw await SampleServiceException.fromResponse(response);
    }

    var bodyJson = await _streamToString(response.stream);
    return factory(jsonDecode(bodyJson));
  }

  Future<ModelStream> _handleResponseStream(http.StreamedResponse response) async {
    if (response.statusCode >= 400) {
      throw await SampleServiceException.fromResponse(response);
    }

    // The http package converts all header names to lower case... so that's a thing.
    return ModelStream(
      content: response.stream,
      contentLength: int.parse(response.headers['content-length'] ?? '0'),
      contentType: response.headers['content-type'] ?? 'application/octet-stream',
      contentFileName: _dispositionFileName(response.headers['content-disposition']),
      contentRange: ModelStreamContentRange.fromString(response.headers['content-range'] ?? ''),
    );
  }

  String _authorize(String callAuthorization) {
    return callAuthorization.trim().isNotEmpty
      ? callAuthorization
      : authorization;
  }

  String _joinUrl(List<String> segments) {
    String stripLeadingSlash(String s) {
      while (s.startsWith('/')) {
        s = s.substring(1);
      }
      return s;
    }
    String stripTrailingSlash(String s) {
      while (s.endsWith('/')) {
        s = s.substring(0, s.length - 1);
      }
      return s;
    }
    bool notEmpty(String s) {
      return s.isNotEmpty;
    }

    return segments
        .map(stripLeadingSlash)
        .map(stripTrailingSlash)
        .where(notEmpty)
        .join('/');
  }

  Map<String, dynamic> _flattenJson(Map<String, dynamic> json) {
    // Adds the given json map entry to the accumulator map. The 'path' contains
    // the period-delimited path for all parent objects we've recurred down from.
    void flattenEntry(String path, String key, dynamic value, Map<String, dynamic> accum) {
      if (value == null) {
        return;
      }

      path = path == '' ? key : '$path.$key';
      if (value is Map<String, dynamic>) {
        value.keys.forEach((key) => flattenEntry(path, key, value[key], accum));
        return;
      }
      accum[path] = value;
    }

    Map<String, dynamic> result = Map<String, dynamic>();
    json.keys.forEach((key) => flattenEntry("", key, json[key], result));
    return result;
  }

  String _dispositionFileName(String? contentDisposition) {
    if (contentDisposition == null) {
      return '';
    }

    var fileNameAttrPos = contentDisposition.indexOf("filename=");
    if (fileNameAttrPos < 0) {
      return '';
    }

    var fileName = contentDisposition.substring(fileNameAttrPos + 9);
    fileName = fileName.startsWith('"') ? fileName.substring(1) : fileName;
    fileName = fileName.endsWith('"') ? fileName.substring(0, fileName.length - 1) : fileName;
    fileName = fileName.replaceAll('\\"', '\"');
    return fileName;
  }
}

class SampleServiceException implements Exception {
  int status;
  String message;

  SampleServiceException(this.status, this.message);

  static Future<SampleServiceException> fromResponse(http.StreamedResponse response) async {
    var body = await _streamToString(response.stream);
    var message = '';
    try {
      Map<String, dynamic> json = jsonDecode(body);
      message = json['message'] ?? json['error'] ?? json['Message'] ?? json['Error'] ?? body;
    }
    catch (_) {
      message = body;
    }
    throw new SampleServiceException(response.statusCode, message);
  }
}


  
  class SampleComplexRequest implements ModelJSON { 
    SampleUser? inUser;
    bool? inFlag;
    double? inFloat;
    TimeTime? inTime;
    TimeTime? inTimePtr;

    SampleComplexRequest({ 
      this.inUser,
      this.inFlag,
      this.inFloat,
      this.inTime,
      this.inTimePtr,
    });

    SampleComplexRequest.fromJson(Map<String, dynamic> json) { 
      
      inUser = SampleUser.fromJson(json['InUser']);
      
      inFlag = json['InFlag'];
      
      inFloat = json['InFloat'];
      
      inTime = json['InTime'];
      
      inTimePtr = json['InTimePtr'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'InUser': inUser?.toJson(),
        
        'InFlag': inFlag,
        
        'InFloat': inFloat,
        
        'InTime': inTime ?? null,
        
        'InTimePtr': inTimePtr ?? null,
      };
    }
  }


  
  class SampleRedirectRequest implements ModelJSON { 

    SampleRedirectRequest();

    SampleRedirectRequest.fromJson(Map<String, dynamic> json) { 
    }

    Map<String, dynamic> toJson() {
      return { 
      };
    }
  }


  
  class SampleRequest implements ModelJSON { 
    String? id;
    String? text;

    SampleRequest({ 
      this.id,
      this.text,
    });

    SampleRequest.fromJson(Map<String, dynamic> json) { 
      
      id = json['ID'];
      
      text = json['Text'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'ID': id,
        
        'Text': text,
      };
    }
  }


  
  typedef MarshalToObject = dynamic;

  
  class SampleDownloadRequest implements ModelJSON { 
    String? format;

    SampleDownloadRequest({ 
      this.format,
    });

    SampleDownloadRequest.fromJson(Map<String, dynamic> json) { 
      
      format = json['Format'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'Format': format,
      };
    }
  }


  
  class SampleResponse implements ModelJSON { 
    String? id;
    String? text;

    SampleResponse({ 
      this.id,
      this.text,
    });

    SampleResponse.fromJson(Map<String, dynamic> json) { 
      
      id = json['ID'];
      
      text = json['Text'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'ID': id,
        
        'Text': text,
      };
    }
  }


  
  typedef CustomDuration = dynamic;

  
  /// SampleUser contains an array of different fields that we support sending to/from clients
  /// in all of our supported languages.
  class SampleUser implements ModelJSON { 
    String? id;
    String? name;
    int? age;
    int? attention;
    CustomDuration? attentionString;
    String? digits;
    MarshalToString? marshalToString;
    MarshalToObject? marshalToObject;

    SampleUser({ 
      this.id,
      this.name,
      this.age,
      this.attention,
      this.attentionString,
      this.digits,
      this.marshalToString,
      this.marshalToObject,
    });

    SampleUser.fromJson(Map<String, dynamic> json) { 
      
      id = json['ID'];
      
      name = json['Name'];
      
      age = json['Age'];
      
      attention = json['Attention'];
      
      attentionString = json['AttentionString'];
      
      digits = json['Digits'];
      
      marshalToString = json['MarshalToString'];
      
      marshalToObject = json['MarshalToObject'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'ID': id,
        
        'Name': name,
        
        'Age': age,
        
        'Attention': attention,
        
        'AttentionString': attentionString,
        
        'Digits': digits,
        
        'MarshalToString': marshalToString ?? null,
        
        'MarshalToObject': marshalToObject ?? null,
      };
    }
  }


  
  typedef TimeDuration = int;

  
  typedef MarshalToString = dynamic;

  
  typedef TimeTime = dynamic;

  
  class SampleComplexResponse implements ModelJSON { 
    bool? outFlag;
    double? outFloat;
    SampleUser? outUser;
    TimeTime? outTime;
    TimeTime? outTimePtr;

    SampleComplexResponse({ 
      this.outFlag,
      this.outFloat,
      this.outUser,
      this.outTime,
      this.outTimePtr,
    });

    SampleComplexResponse.fromJson(Map<String, dynamic> json) { 
      
      outFlag = json['OutFlag'];
      
      outFloat = json['OutFloat'];
      
      outUser = SampleUser.fromJson(json['OutUser']);
      
      outTime = json['OutTime'];
      
      outTimePtr = json['OutTimePtr'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'OutFlag': outFlag,
        
        'OutFloat': outFloat,
        
        'OutUser': outUser?.toJson(),
        
        'OutTime': outTime ?? null,
        
        'OutTimePtr': outTimePtr ?? null,
      };
    }
  }


  
  class SampleDownloadResponse extends ModelStream {
    SampleDownloadResponse({
      Stream<List<int>>? content,
      int? contentLength,
      String? contentType,
      String? contentFileName,
      ModelStreamContentRange? contentRange,
    }) : super(content: content, contentType: contentType, contentFileName: contentFileName, contentLength: contentLength, contentRange: contentRange);
  }

  
  class SampleSecurityResponse implements ModelJSON { 
    List<String>? roles;

    SampleSecurityResponse({ 
      this.roles,
    });

    SampleSecurityResponse.fromJson(Map<String, dynamic> json) { 
      
      roles = _map(json['Roles'], (x) => x);
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'Roles': _map(roles, (x) => x),
      };
    }
  }


  
  class SampleSecurityRequest implements ModelJSON { 
    String? id;
    SampleUser? user;

    SampleSecurityRequest({ 
      this.id,
      this.user,
    });

    SampleSecurityRequest.fromJson(Map<String, dynamic> json) { 
      
      id = json['ID'];
      
      user = SampleUser.fromJson(json['User']);
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'ID': id,
        
        'User': user?.toJson(),
      };
    }
  }


  
  class SampleRedirectResponse extends ModelStream {
    SampleRedirectResponse({
      Stream<List<int>>? content,
      int? contentLength,
      String? contentType,
      String? contentFileName,
      ModelStreamContentRange? contentRange,
    }) : super(content: content, contentType: contentType, contentFileName: contentFileName, contentLength: contentLength, contentRange: contentRange);
  }


class ModelJSON {
  Map<String, dynamic> toJson() {
    throw new Exception('toJson not implemented');
  }
}

class ModelStream implements ModelJSON {
  Stream<List<int>>? content;
  int? contentLength;
  String? contentType;
  String? contentFileName;
  ModelStreamContentRange? contentRange;

  ModelStream({
    this.content, this.contentLength, this.contentType, this.contentFileName, this.contentRange,
  });

  Map<String, dynamic> toJson() {
    return {
      'Content': _streamToString(content),
      'ContentLength': contentLength ?? 0,
      'ContentType': contentType ?? 'application/octet-stream',
      'ContentFileName': contentFileName ?? '',
      'ContentRange': (contentRange ?? ModelStreamContentRange()).toJson(),
    };
  }
}

class ModelStreamContentRange {
  String unit;
  int start;
  int end;
  int size;

  ModelStreamContentRange({this.unit = 'bytes', this.start = 0, this.end = 0, this.size = 0});

  static ModelStreamContentRange fromString(String rangeHeader) {
    rangeHeader = rangeHeader.trim();
    if (rangeHeader == '') {
      return ModelStreamContentRange();
    }

    var tokens = rangeHeader.split(RegExp('[ -/]'));
    switch (tokens.length) {
    case 3: return ModelStreamContentRange(
      start: int.parse(tokens[0]),
      end: int.parse(tokens[1]),
      size: int.parse(tokens[2]),
    );
    case 4: return ModelStreamContentRange(
      unit: tokens[0],
      start: int.parse(tokens[1]),
      end: int.parse(tokens[2]),
      size: int.parse(tokens[3]),
    );
    default:
      // Invalid header format so just pretend like you didn't get anything
      return ModelStreamContentRange();
    }
  }

  Map<String, dynamic> toJson() {
    return { 'Unit': unit, 'Start': start, 'End': end, 'Size': size };
  }
}

List<T>? _map<T>(List<dynamic>? jsonList, T Function(dynamic) mapping) {
  return jsonList == null ? null : jsonList.map(mapping).toList();
}

Future<String> _streamToString(Stream<List<int>>? stream) async {
  if (stream == null) {
    return '';
  }
  var bodyCompleter = new Completer<String>();
  stream.transform(utf8.decoder).listen(bodyCompleter.complete);
  return bodyCompleter.future;
}








