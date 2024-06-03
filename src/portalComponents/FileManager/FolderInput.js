import React, { useEffect, useRef } from 'react';

export default function FolderInput(props) {
  const { onChange } = props;
  const ref = useRef(null);

  useEffect(() => {
    if (ref.current) {
      ref.current.setAttribute("directory", "");
      ref.current.setAttribute("webkitdirectory", "");
      ref.current.setAttribute("mozdirectory", "");
    }
  }, [ref]);

  return (
    <input
      type="file"
      ref={ref}
      onChange={onChange}
    />
  )
}