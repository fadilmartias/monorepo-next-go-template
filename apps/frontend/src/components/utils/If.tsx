import React from "react";

export const Then = ({ children }: { children: React.ReactNode }) => <>{children}</>;
export const Else = ({ children }: { children: React.ReactNode }) => <>{children}</>;

export const If = ({
  condition,
  children,
}: {
  condition: boolean;
  children: React.ReactNode;
}) => {
  let thenChild = null;
  let elseChild = null;

  React.Children.forEach(children, (child) => {
    if (!React.isValidElement(child)) return;
    if (child.type === Then) thenChild = child;
    if (child.type === Else) elseChild = child;
  });

  return <>{condition ? thenChild : elseChild}</>;
};
