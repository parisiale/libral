#pragma once

#include <vector>
#include <memory>
#include <boost/optional.hpp>

#include <yaml-cpp/yaml.h>

#include <libral/result.hpp>
#include <libral/type.hpp>
#include <libral/provider.hpp>
#include <libral/environment.hpp>

namespace libral {
  class ral : public std::enable_shared_from_this<ral> {
  public:
    std::vector<std::unique_ptr<type>> types(void);
    boost::optional<std::unique_ptr<type>> find_type(const std::string& name);

    /* Create an instance of the RAL */
    static std::shared_ptr<ral> create(std::vector<std::string> data_dirs);

    const std::vector<std::string>& data_dirs() const { return _data_dirs; }

    boost::optional<std::string>
    find_in_data_dirs(const std::string& file) const;

  protected:
    ral(const std::vector<std::string>& data_dir);

    environment make_env() { return environment(shared_from_this()); }
  private:
    bool add_type(std::vector<std::unique_ptr<type>>& types,
                  environment &env,
                  const std::string& name, std::shared_ptr<provider> prov);
    result<YAML::Node> parse_metadata(const std::string& path,
                                      const std::string& yaml) const;
    result<std::string> run_describe(const std::string& path) const;
    std::vector<std::string> _data_dirs;
  };
}
