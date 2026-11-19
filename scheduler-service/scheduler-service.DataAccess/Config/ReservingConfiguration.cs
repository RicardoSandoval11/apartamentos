using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using scheduler_service.DataAccess.Entities;

namespace scheduler_service.DataAccess.Config;

public class ReservingConfiguration : IEntityTypeConfiguration<Reserving>
{
    public void Configure(EntityTypeBuilder<Reserving> builder)
    {
        builder.ToTable("reserving");

        builder.HasKey(e => e.Id);

        builder.Property(e => e.Id)
        .ValueGeneratedOnAdd()
        .HasColumnName("id")
        .IsRequired();

        builder.Property(e => e.Name)
        .HasColumnName("name")
        .IsRequired();

        builder.Property(e => e.Email)
        .HasColumnName("email")
        .IsRequired();

        builder.Property(e => e.PhoneNumber)
        .HasColumnName("phone_number")
        .IsRequired();

        builder.HasOne(e => e.Reservation)
        .WithOne(e => e.Reserving);
    }
}