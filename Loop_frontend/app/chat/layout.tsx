export default function ChatPageLayout({
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
