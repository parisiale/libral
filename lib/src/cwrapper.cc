#include <libral/cwrapper.hpp>

#include <libral/ral.hpp>
#include <libral/emitter/json_emitter.hpp>

#include <string>
#include <vector>

namespace lib = libral;

/* TODO in GOLand:
 *
 * getTypes()(string, error) --> string == MUST serialize vector as JSON array
 * getAllResources(type) (string, error) --> string == JSON object serialized
 * getResource(type, resource) (string, error) --> string == JSON object serialized
 * getSchema(type)(string, error) --> how to derive a JSON schema??? Use JsonSchema???
 *
 */

//
// Helpers
//

static char* str_to_cstr(const std::string& s) {
    char * cs = new char [s.length()+1];
    std::strcpy (cs, s.c_str());
    return cs;
}

//
// Public interface
//

// TODO(ale): decide between pointers & structs for calling / returning
//            conventions and use a single approach

char* get_all(char* type_name)
{
    std::vector<std::string> data_dirs;
    auto ral = lib::ral::create(data_dirs);
    auto opt_type = ral->find_type(std::string(type_name));

    if (!opt_type)
        return NULL;

    auto results = (*opt_type)->instances();
    lib::json_emitter em {};
    auto resources = em.parse_list(**opt_type, results);
    return str_to_cstr(resources);
}

// TODO(ale): make the pointer approach work
uint8_t get_all_with_err(char **resource,
                         char *type_name) {
    char *tmp = get_all(type_name);

    if (resource == nullptr)
        return 1;

    resource = &tmp;
    return 0;
}

struct outcome get_all_outcome(char* type_name_c) {
    char *tmp = get_all(type_name_c);
    uint8_t e_c = (tmp != NULL) ? 0 : 1;
    struct outcome out = {tmp, e_c};
    return out;
}

uint8_t get_types(char *types) {
    return 0;
}
