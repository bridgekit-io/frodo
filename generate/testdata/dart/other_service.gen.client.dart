// Code generated by Frodo - DO NOT EDIT.
//
//   Timestamp: Mon, 21 Oct 2024 14:27:35 EDT
//   Source:    other_service.go
//   Generator: https://github.com/bridgekitio/frodo
//
import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;

class OtherServiceClient {
  final String baseURL;
  String authorization;
  http.Client httpClient = http.Client();

  OtherServiceClient(this.baseURL, {
      this.authorization = '',
  });

  
  /// ChainFail fires after ChainOne, but should always return an error. This will prevent ChainFailAfter
  /// from ever actually running.
  Future<OtherResponse> ChainFail(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainFail';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ChainFailAfter is dependent on a successful call to ChainFail... which always fails. So this NEVER runs.
  Future<OtherResponse> ChainFailAfter(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainFailAfter';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ChainFour is used to test that methods invoked via the event gateway can trigger even more events.
  Future<OtherResponse> ChainFour(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainFour';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ChainOne allows us to test the cascading of events to create more complex flows. When this
  /// finishes it will trigger ChainTwo which will, in turn, trigger ChainThree and ChainFour.
  Future<OtherResponse> ChainOne(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainOne';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ChainThree is used to test that methods invoked via the event gateway can trigger even more events.
  Future<OtherResponse> ChainThree(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainThree';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ChainTwo is used to test that methods invoked via the event gateway can trigger even more events.
  Future<OtherResponse> ChainTwo(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ChainTwo';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// ListenWell can listen for successful responses across multiple services.
  Future<OtherResponse> ListenWell(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.ListenWell';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// RPCExample invokes the TriggerUpperCase() function on the SampleService to get work done.
  /// This will make sure that we can do cross-service communication.
  Future<OtherResponse> RPCExample(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.RPCExample';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
    }
  }
  
  /// SpaceOut takes your input text and puts spaces in between all the letters.
  Future<OtherResponse> SpaceOut(OtherRequest serviceRequest, {String authorization = ''}) async {
  
    var requestJson = serviceRequest.toJson();
    var method = 'POST';
    var route = '/OtherService.SpaceOut';
    var uri = _joinUrl([baseURL, _buildRequestPath(method, route, requestJson)]);

    try {
      final request = http.Request(method, Uri.parse(uri));
      request.headers['Accept'] = 'application/json';
      request.headers['Authorization'] = _authorize(authorization);
      request.headers['Content-Type'] = 'application/json';
      request.body = jsonEncode(requestJson);

      final response = await httpClient.send(request);
      return _handleResponseJson(response, (json) => OtherResponse.fromJson(json));
    } on OtherServiceException catch (e) {
      throw e; // already has status information
    } catch (e) {
      throw OtherServiceException(500, e.toString());
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
      throw await OtherServiceException.fromResponse(response);
    }

    var bodyJson = await _streamToString(response.stream);
    return factory(jsonDecode(bodyJson));
  }

  Future<ModelStream> _handleResponseStream(http.StreamedResponse response) async {
    if (response.statusCode >= 400) {
      throw await OtherServiceException.fromResponse(response);
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

class OtherServiceException implements Exception {
  int status;
  String message;

  OtherServiceException(this.status, this.message);

  static Future<OtherServiceException> fromResponse(http.StreamedResponse response) async {
    var body = await _streamToString(response.stream);
    var message = '';
    try {
      Map<String, dynamic> json = jsonDecode(body);
      message = json['message'] ?? json['error'] ?? json['Message'] ?? json['Error'] ?? body;
    }
    catch (_) {
      message = body;
    }
    throw new OtherServiceException(response.statusCode, message);
  }
}


  
  /// OtherRequest is a basic payload that partially matches the schema of SampleResponse so
  /// when we invoke service methods through the event gateway, we can make sure that we
  /// can get the Text value while ignoring everything else from the original payload.
  class OtherRequest implements ModelJSON { 
    bool? uniqueThing;
    String? text;

    OtherRequest({ 
      this.uniqueThing,
      this.text,
    });

    OtherRequest.fromJson(Map<String, dynamic> json) { 
      
      uniqueThing = json['UniqueThing'];
      
      text = json['Text'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'UniqueThing': uniqueThing,
        
        'Text': text,
      };
    }
  }


  
  class OtherResponse implements ModelJSON { 
    bool? uniqueThing;
    String? text;

    OtherResponse({ 
      this.uniqueThing,
      this.text,
    });

    OtherResponse.fromJson(Map<String, dynamic> json) { 
      
      uniqueThing = json['UniqueThing'];
      
      text = json['Text'];
    }

    Map<String, dynamic> toJson() {
      return { 
        
        'UniqueThing': uniqueThing,
        
        'Text': text,
      };
    }
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








