export default function CreatePageLayout({
    children,
  }: {
    children: React.ReactNode;
  }) {
    return (
      <section className="flex flex-col items-left justify-left gap-4 py-8 md:py-10">
        <div className="inline-block text-left justify-left">
          {children}
        </div>
      </section>
    );
  }
  