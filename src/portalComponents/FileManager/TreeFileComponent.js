import {
  UncontrolledTreeEnvironment,
  Tree,
  StaticTreeDataProvider,
  TreeItemIndex,
  TreeItem
} from "react-complex-tree";
import React from "react";
import "react-complex-tree/lib/style-modern.css";
import { renderers as bpRenderers } from 'react-complex-tree-blueprintjs-renderers';

export default function TreeFileComponent(props) {
  const { data, root, onSelectItems } = props;

  return (
    <UncontrolledTreeEnvironment
    {...bpRenderers}
      dataProvider={
        new StaticTreeDataProvider(data, (item, data) => ({
          ...item,
          data
        }))
      }
      getItemTitle={(item) => {
        const splitPath = item.data.split("/");
        return splitPath[splitPath.length - 1];
      }}
      viewState={{}}
      onSelectItems={onSelectItems}
      canDragAndDrop={true}
      canDropOnFolder={true}
    >
      <Tree treeId={root} rootItem={root} treeLabel={root} />
    </UncontrolledTreeEnvironment>
  );
}