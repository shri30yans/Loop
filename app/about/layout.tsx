export default function AboutLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section className="flex flex-col items-center justify-left gap-4 py-8 md:py-10">
      <div className="inline-block max-w-screen text-left justify-left text-xl">
        {children}
      </div>
    </section>
  );
}
