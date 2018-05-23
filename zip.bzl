def _pkg_zip_impl(ctx):
  files = " ".join(["%s" % dep.path for dep in ctx.files.deps])
  zipper_arg_path = ctx.actions.declare_file("%s_zipper_args" % ctx.outputs.zip.path)
  #label = dep.label
  #query = "$(locations %s)" % label
  #location = ctx.expand_location(query, [dep]).split(" ")[0]
  
  #ctx.file_action(zipper_arg_path, files)
  ctx.actions.write(zipper_arg_path, files)
  print(["the files we got are", files, "the zipper_arg_path is", zipper_arg_path.path])
  cmd = """
rm -f {out}
{zipper} cf {out} @{path}
"""
  cmd = cmd.format(out=ctx.outputs.zip.path,
                   path=zipper_arg_path.path,
                   zipper=ctx.executable._zipper.path)

  ctx.actions.run_shell(
      inputs = ctx.files.deps + [ctx.executable._zipper, zipper_arg_path],
      outputs = [ctx.outputs.zip],
      command = cmd,
      progress_message = "making pkg_zip archive %s (%s files)" % (ctx.label, len(files)),
  )

pkg_zip = rule(
    implementation=_pkg_zip_impl,
    attrs={
        "deps": attr.label_list(allow_files = True),
        "_zipper": attr.label(
            default=Label("@bazel_tools//tools/zip:zipper"), 
            allow_files=True,
            executable=True, 
            cfg="host",
        ),
    },
    outputs={
        "zip": "pkg_%{name}"
    },
)