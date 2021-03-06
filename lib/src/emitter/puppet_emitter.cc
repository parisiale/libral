#include <libral/emitter/puppet_emitter.hpp>

#include <leatherman/logging/logging.hpp>

#include <iostream>
#include <iomanip>

using namespace leatherman::logging;
using namespace leatherman::locale;

namespace libral {

  namespace color {
  // Very poor man's output coloring. Call init() to fill these
  // with color escape sequences
  std::string cyan, green, yellow, red, magenta, blue, reset;

  void init() {
    if (isatty(1)) {
      cyan = "\33[0;36m";
      green = "\33[0;32m";
      yellow = "\33[0;33m";
      red = "\33[0;31m";
      magenta = "\33[0;35m";
      blue = "\33[0;34m";

      reset = "\33[0m";
    }
  } } // namespace color


  puppet_emitter::puppet_emitter() {
    // Color handling is horrible: as a side-effect of crating a puppet_emitter,
    // we initialize the global color constants.
    color::init();
  }

  void puppet_emitter::print_set(const type &type,
                       const result<std::pair<update, changes>>& rslt) {
    if (rslt) {
      const auto& upd = rslt->first;
      const auto &chgs = rslt->second;
      print_resource(type, type.prov().create(upd, chgs));
      std::cout << chgs << std::endl;
    } else {
      std::cout << color::red << _("failed: {1}", rslt.err().detail) << color::reset << std::endl;
    }
  }

  void puppet_emitter::print_find(const type &type,
                       const result<boost::optional<resource>> &inst) {
    if (!inst) {
      std::cerr << color::red <<
        _("failed: {1}", inst.err().detail) << color::reset << std::endl;
      return;
    }
    if (inst.ok()) {
      print_resource(type, *inst.ok());
    }
  }

  void puppet_emitter::print_list(const type &type,
                       const result<std::vector<resource>>& rslt) {
    if (!rslt) {
      std::cout << color::red << _("failed: {1}", rslt.err().detail) << color::reset << std::endl;
      return;
    }
    for (const auto& inst : rslt.ok()) {
      print_resource(type, inst);
    }
  }

  void puppet_emitter::print_types(const std::vector<std::unique_ptr<type>>& types) {
    for (const auto& t : types) {
      std::cout << t->qname() << std::endl;
    }
  }

  void puppet_emitter::print_resource(const type &type, const resource &res) {
    std::cout << color::blue << type.qname() << color::reset << " { '"
         << color::blue << res.name() << color::reset << "':" << std::endl;
    uint16_t maxlen = 0;
    for (const auto& a : res.attrs()) {
      if (a.first.length() > maxlen) maxlen = a.first.length();
    }
    auto ens = res.lookup<std::string>("ensure");
    if (ens) {
      print_resource_attr("ensure", *ens, maxlen);
    }
    for (const auto& a : res.attrs()) {
      if (a.first == "ensure")
        continue;
      print_resource_attr(a.first, a.second, maxlen);
    }
    std::cout << "}" << std::endl;
  }

  void puppet_emitter::print_resource_attr(const std::string& name,
                                           const value& v,
                                           uint16_t maxlen) {
    std::cout << "  " << color::green << std::left << std::setw(maxlen)
              << name << color::reset
              << " => " << v << "," << std::endl;
  }
}
