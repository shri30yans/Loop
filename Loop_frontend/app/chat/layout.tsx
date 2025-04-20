export default function CreatePageLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section className="h-full">
      {children}
    </section>
  );
}
