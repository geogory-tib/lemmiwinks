#include <cjson/cJSON.h>
#include "../include/gtd.h"
const char test_json[] =
"{\n"
"  \"name\": \"Test Parser\",\n"
"  \"version\": 1.0,\n"
"  \"valid\": true,\n"
"  \"notes\": \"Escapes: \\\"quote\\\", \\\\slash\\\\, unicode: \\u263A\",\n"
"  \"values\": [1, 2.5, -3, 4e2, null, false, true],\n"
"  \"nested\": {\n"
"    \"array\": [\n"
"      { \"id\": 1, \"tag\": \"alpha\" },\n"
"      { \"id\": 2, \"tag\": \"beta\", \"extra\": [\"x\", \"y\", \"z\"] }\n"
"    ],\n"
"    \"deep\": {\n"
"      \"flag\": false,\n"
"      \"number\": 123456789,\n"
"      \"empty_obj\": {},\n"
"      \"empty_array\": []\n"
"    }\n"
"  }\n"
"}";

int main(int argc, char **argv) {
  auto object = cJSON_ParseWithLength((char *)test_json, sizeof(test_json));
  auto item = cJSON_GetObjectItemCaseSensitive(object, "name");
  if (!item) {
	panic("Item is null");
  }
  if (cJSON_IsString(item)) {
	printf("%s\n",item->valuestring);
  }
  return 0;
}
