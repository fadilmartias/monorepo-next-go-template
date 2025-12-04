import Typography from "./typography";

interface Props {
  title: string;
  desc?: string;
}

export default function HomepageTitle({ title, desc }: Props) {
  return (
    <div>
      <div className="flex items-center gap-2">
        <Typography variant="h5" as="h2" className="text-primary">
          {title}
        </Typography>
      </div>
      <Typography variant="b3" dangerouslySetInnerHTML={{ __html: desc || "Ayo pesan sekarang sebelum kehabisan!" }} />
    </div>
  );
}
